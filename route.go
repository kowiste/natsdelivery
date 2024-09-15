package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func sendData(c *gin.Context) {
	var request struct {
		Topic string `json:"topic"`
		Event string `json:"event"`
		Data  any    `json:"data"`
	}

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = nd.Send(request.Topic, request.Event, request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Message sent successfully"})
}
