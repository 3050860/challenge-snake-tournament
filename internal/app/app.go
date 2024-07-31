package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"snake-tournament/internal/config"
	"snake-tournament/internal/tournament/database"
	"snake-tournament/internal/tournament/handler"
	"snake-tournament/internal/tournament/service"
	"snake-tournament/pkg/logging"
	"snake-tournament/pkg/metrics"
	"snake-tournament/pkg/mongodb_client"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("router initializing")
	router := httprouter.New()

	logger.Println("swagger docs initialization")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Println("heartbeat metric initializing")
	metricHandler := metrics.Handler{}
	metricHandler.Register(router)

	logger.Println("mongodb connection initializing")
	mongodbClient, err := mongodb_client.NewClient(context.Background(), cfg.MongoDB.ConnectString, cfg.MongoDB.Database, cfg.MongoDB.AuthDatabase, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		panic(err)
	}

	logger.Println("ticket service initializing")
	ticketService := service.NewTicketService(cfg)

	logger.Println("games database initializing")
	gamesDatabase := database.NewGamesDatabase(mongodbClient, cfg.MongoDB.Collection)

	logger.Println("game service initializing")
	gameService := service.NewGameService(*gamesDatabase, ticketService)

	logger.Println("user service initializing")
	userService := service.NewUserService(cfg)

	logger.Println("game handler initializing")
	handler.NewRecordHandler(gameService, router, userService)

	return App{
		cfg,
		logger,
		router,
		nil,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == config.ListenTypeSock {
		appDir, err := filepath.Abs(os.Args[0])
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"https://localhost:3000", "https://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Authorization", "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		Debug:              false,
	})

	a.httpServer = &http.Server{
		Handler:      c.Handler(a.router),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely initialized and started")
	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}

}
