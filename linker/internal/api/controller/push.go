package controller

import (
	"net/http"
	"push/common/lib"
	"push/linker/internal/api/dto"
	"push/linker/internal/api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PushController struct {
	logger  lib.Logger
	service service.PushService
}

func NewPushController(logger lib.Logger, service service.PushService) PushController {
	return PushController{
		service: service,
		logger:  logger,
	}
}

func (p PushController) PostPush(c *gin.Context) {
	userIdStr := c.Param("userid")
	uid, err := strconv.Atoi(userIdStr)
	if err != nil || uid < 0 {
		c.JSON(400, gin.H{"error": "Invalid userid"})
		return
	}
	userId := uint(uid)

	var req dto.CreateMessageReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dto.CreateMessageDTO{
		UserID:  userId,
		Title:   req.Title,
		Content: req.Content,
	}

	ctx := c.Request.Context()
	id, err := p.service.PostPush(ctx, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error message"}) // TODO : 에러 처리 개선 필요.
		return
	}
	c.JSON(200, gin.H{"message_id": id})
}
