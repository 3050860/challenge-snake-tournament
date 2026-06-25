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
	"snake-tournament/internal/clients/mongodb"
	"snake-tournament/internal/config"
	"snake-tournament/internal/handlers"
	"snake-tournament/internal/repository"
	"snake-tournament/internal/services"
	"snake-tournament/pkg/metrics"
	"time"

	"github.com/alexsuslov/ehttp/pkg/hticket"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg        *config.Config
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(cfg *config.Config) (App, error) {
	router := httprouter.New()
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)
	metricHandler := metrics.Handler{}
	metricHandler.Register(router)

	logrus.Info("connecting to mongodb")
	mongodbClient, err := mongodb.NewClient(context.Background(), cfg.MongoDB.ConnectString, cfg.MongoDB.Database, cfg.MongoDB.AuthDatabase, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		panic(err)
	}

	ticketService := hticket.NewHService(cfg.HTicketConfig)
	gamesRepo := repository.NewGames(mongodbClient, cfg.MongoDB.Collection)
	gameService := services.NewGameService(gamesRepo, ticketService)

	handlers.NewRecordHandler(gameService, router)

	return App{
		cfg,
		router,
		nil,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	var listener net.Listener

	if a.cfg.Listen.Type == config.ListenTypeSock {
		appDir, err := filepath.Abs(os.Args[0])
		if err != nil {
			logrus.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		logrus.Infof("socket path: %s", socketPath)
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			logrus.Fatal(err)
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

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logrus.Warn("server shutdown")
		default:
			logrus.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

}
