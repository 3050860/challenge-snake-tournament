package config

import (
	"github.com/alexsuslov/ehttp/pkg/hticket"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		SocketFile string `env:"SOCKET_FILE" env-default:"app.sock"`
		Type       string `env:"LISTEN_TYPE" env-default:"port"`
		BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"PORT" env-default:"10003"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL" env-default:"trace"`
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-default:"admin"`
			Password string `env:"ADMIN_PWD" env-default:"admin"`
		}
	}
	MongoDB struct {
		ConnectString string `env:"MONGO_URL" env-default:"localhost"`
		Database      string `env:"DATABASE" env-default:"snake_tournament"`
		Collection    string `env:"GAMES_COLLECTION" env-default:"games"`
		AuthDatabase  string `env:"AUTH_COLLECTION" env-default:"admin"`
		Username      string `env:"MONGODB_USERNAME" env-default:"root"`
		Password      string `env:"MONGODB_PASSWORD" env-default:"root"`
	}
	UserService struct {
		Host string `env:"USERS_SERVICE_HOST" env-default:"localhost:10002"`
	}
	HTicketConfig hticket.Config
	TicketService struct {
		Host string `env:"TICKETS_SERVICE_HOST" env-default:"localhost:10004"`
	}
	Keys struct {
		AccessKey  string `env:"ACCESS_KEY" env-default:"18d8debd1eec2eb338c4a9a8815633cede19cf3d17b0f20c60cf3839a89699cb"`
		JWTSignKey string `env:"JWT_SIGN_KEY" env-default:"alsfjak12h4i1h2uas7f7241231o1u2io5u12asopua0w9812"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "An error occurred during reading config"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
