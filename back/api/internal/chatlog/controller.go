package chatlog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatLog struct {
	ClientID string `json:"client_id"`
	Message  string `json:"message"`
}

func PostChatLog(c *gin.Context) {
	var chatLog ChatLog
	if err := c.ShouldBindJSON(&chatLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[CHAT LOG] %s: %s\n", chatLog.ClientID, chatLog.Message)

	c.JSON(http.StatusOK, gin.H{"status": "logged"})
}
