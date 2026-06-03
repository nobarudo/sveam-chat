package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type ConfigList struct {
	Port   string
	DBName string
}

var config ConfigList

func init() {
	LoadConfig()
}

func LoadConfig() {
	if err := godotenv.Load("configs/server.env"); err != nil {
		slog.Info("環境変数ファイルの読み込みスキップ (システム環境変数を使用)", "err", err)
	}
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	config.DBName = os.Getenv("DB_NAME")
	if config.DBName == "" {
		config.DBName = "hemochart.db"
	}
}

func GetConfig() ConfigList {
	return config
}
