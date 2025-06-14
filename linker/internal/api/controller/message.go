package controller

import (
	"net/http"
	commondto "push/common/dto"
	"push/common/lib"
	"push/linker/dto"
	"push/linker/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	logger  lib.Logger
	service service.MessageService
}

func NewMessageController(logger lib.Logger, service service.MessageService) MessageController {
	return MessageController{
		service: service,
		logger:  logger,
	}
}

func (p MessageController) UpdateStatusToReceive(c *gin.Context) {
	msgIdParam := c.Param("msgid")
	msgId, err := strconv.Atoi(msgIdParam)
	if err != nil || msgId < 0 {
		c.JSON(400, gin.H{"error": "Invalid msgid"})
		return
	}
	msgIdUint64 := uint64(msgId)

	var req dto.UpdateStatusReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	if err := p.service.ReceiveMessage(ctx, msgIdUint64); err != nil {
		p.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error message"}) // TODO : 에러 처리 개선 필요.
		return
	}

	res := commondto.CommonResponse[dto.UpdateStatusRes]{}
	res.Data = dto.UpdateStatusRes{MsgId: msgIdUint64}
	c.JSON(200, res)

}
