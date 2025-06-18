package controller

import (
	"net/http"
	commondto "push/common/dto"
	"push/common/lib/logger"
	re "push/linker/dto"
	"push/linker/internal/api/dto"
	"push/linker/internal/service"
	"push/linker/internal/worker"
	"push/linker/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	logger          *logger.Logger
	service         service.MessageService
	jobUpdateStatus *worker.JobUpdateStatus
}

func NewMessageController(logger *logger.Logger, service service.MessageService, jobUpdateStatus *worker.JobUpdateStatus) MessageController {
	return MessageController{
		service: service,
		logger:  logger,

		jobUpdateStatus: jobUpdateStatus,
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

	var reqDto re.UpdateStatusReq

	if err := c.ShouldBindJSON(&reqDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ctx := c.Request.Context()
	dto := dto.UpdateMessageDTO{
		Id:     msgIdUint64,
		Status: types.StatusSent,
	}
	if err := p.jobUpdateStatus.Enqueue(dto); err != nil {
		p.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error message"}) // TODO : 에러 처리 개선 필요.
		return
	}

	res := commondto.CommonResponse[re.UpdateStatusRes]{}
	res.Data = re.UpdateStatusRes{MsgId: msgIdUint64}
	c.JSON(200, res)

}
