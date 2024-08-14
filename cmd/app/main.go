package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	_ "snake-tournament/docs"
	"snake-tournament/internal/app"
	"snake-tournament/internal/config"
)

func main() {
	logrus.Info("Starting Snake Tournament ...")
	if err := godotenv.Load(); err != nil {
		logrus.Warning("no .env file")
	}
	cfg := config.GetConfig()
	a, err := app.NewApp(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	a.Run()
}
