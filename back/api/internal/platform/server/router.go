package server

import (
	"svem-chat-api/internal/chatlog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func routing(router *gin.Engine, db *gorm.DB) {
	ChatlogRepo := chatlog.NewChatlogRepository(db)
	ChatlogController := chatlog.NewChatlogController(ChatlogRepo)

	router.POST("/", ChatlogController.PostChatLog)
}
