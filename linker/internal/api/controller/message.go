package controller

import (
	"push/common/lib"
	"push/linker/internal/api/service"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	logger  lib.Logger
	service service.PushService
}

func NewMessageController(logger lib.Logger, service service.PushService) MessageController {
	return MessageController{
		service: service,
		logger:  logger,
	}
}

func (m MessageController) Test(c *gin.Context) {
	m.service.Test(c)
	c.JSON(200, gin.H{"data": "success"})
}
