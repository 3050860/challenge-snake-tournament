package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	_ "snake-tournament/docs"
	"snake-tournament/internal/app"
	"snake-tournament/internal/config"
	"snake-tournament/pkg/logging"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warning("no .env file")
	}
	log.Print("config initialization")
	cfg := config.GetConfig()

	log.Printf("logging initialized.")
	logger := logging.InitLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("running application")
	a.Run()
}
