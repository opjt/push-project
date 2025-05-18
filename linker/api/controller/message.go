package controller

import (
	"push/common/lib"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	logger lib.Logger
	// service service.WebhookService
}

func NewMessageController(logger lib.Logger) MessageController {
	return MessageController{
		// service: userService,
		logger: logger,
	}
}

func (w MessageController) Test(c *gin.Context) {

	c.JSON(200, gin.H{"data": "success"})
}
