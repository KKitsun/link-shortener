package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env string
	HTTPServer
	DBConfig
}

type HTTPServer struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func MustLoad() *Config {
	if err := godotenv.Load("local.env"); err != nil {
		log.Fatalf("error loading env variables: %s", err)
	}

	viper.AddConfigPath(os.Getenv("CONFIG_PATH"))
	viper.SetConfigName(os.Getenv("CONFIG_NAME"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err)
	}

	cfg := Config{
		Env: viper.GetString("env"),
		HTTPServer: HTTPServer{
			Port: viper.GetString("http_server.port"),
		},
		DBConfig: DBConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBname:   viper.GetString("db.dbname"),
		},
	}

	return &cfg
}
