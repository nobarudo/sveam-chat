package database

import (
	"log/slog"

	"svem-chat-api/internal/chatlog"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(&chatlog.ChatLog{})
	slog.Info("db migrate")
}
