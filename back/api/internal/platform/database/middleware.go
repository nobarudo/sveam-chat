package database

import (
	"log"

	"svem-chat-api/internal/platform/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg *config.ConfigList) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("DB connected:", cfg.DBName)

	return db, nil
}
