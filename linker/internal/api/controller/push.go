package controller

import (
	"net/http"
	"push/common/lib/logger"
	"push/linker/internal/api/dto"
	"push/linker/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PushController struct {
	logger  *logger.Logger
	service service.PushService
}

func NewPushController(logger *logger.Logger, service service.PushService) PushController {
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
	userId := uint64(uid)

	var req dto.PostPushReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := dto.PostPushDTO{
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
	}

	ctx := c.Request.Context()
	id, err := p.service.PostPush(ctx, dto)
	if err != nil {
		p.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error message"}) // TODO : 에러 처리 개선 필요.
		return
	}
	c.JSON(200, gin.H{"message_id": id})
}
