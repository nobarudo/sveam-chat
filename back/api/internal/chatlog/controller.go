package chatlog

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatlogController struct {
	repo *ChatlogRepository
}

func NewChatlogController(repo *ChatlogRepository) *ChatlogController {
	return &ChatlogController{repo: repo}
}

type ChatLog struct {
	ClientID string `json:"client_id"`
	Message  string `json:"message"`
}

func (cc *ChatlogController) PostChatLog(c *gin.Context) {
	var chatLog ChatLog
	if err := c.ShouldBindJSON(&chatLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("[CHAT LOG]", "userid:", chatLog.ClientID, "message:", chatLog.Message)
	err := cc.repo.createChatlog(chatLog)
	if err != nil {
		slog.Error("データベースエラー", err)
	}

	c.JSON(http.StatusOK, gin.H{"status": "logged"})
}
