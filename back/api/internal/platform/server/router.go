package server

import (
	"svem-chat-api/internal/chatlog"

	"github.com/gin-gonic/gin"
)

func routing(router *gin.Engine) {
	router.POST("/", chatlog.PostChatLog)
}
