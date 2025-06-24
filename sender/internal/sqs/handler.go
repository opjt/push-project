package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"push/common/lib/logger"
	msgTypes "push/linker/types"
	"push/sender/internal/dto"
	"push/sender/internal/service"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler interface {
	HandleMessage(ctx context.Context, msg types.Message) error
}

type handler struct {
	log           *logger.Logger
	senderService service.SenderService
	// mclient       client.MessageClient
	// sessionClient sclient.SessionClient
}

func NewHandler(log *logger.Logger, senderService service.SenderService) Handler {
	return &handler{
		log:           log,
		senderService: senderService,
	}
}

func (h *handler) HandleMessage(ctx context.Context, msg types.Message) error {
	pushMsg, err := parseSqsMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	h.log.Infof("Received push message: %+v", pushMsg)
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	updateStatusReq := dto.UpdateMessageStatus{
		Id:       uint64(pushMsg.MsgID),
		Status:   msgTypes.StatusSending,
		SnsMsgId: *msg.MessageId,
	}
	// Linker에게 MessageStatus Update 요청
	if err = h.senderService.UpdateMessage(ctx, updateStatusReq); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}
	// SessionManager에게 메세지 전송 요청
	return h.sendPushMessage(pushMsg)
}

func (h *handler) sendPushMessage(pushMsg *dto.PushMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := h.senderService.PushMessage(ctx, *pushMsg); err != nil {
		return fmt.Errorf("failed to send push message: %w", err)
	}
	return nil
}

// SQS.Message 를 json.Unmarshal 하는 함수
func parseSqsMessage(msg types.Message) (*dto.PushMessage, error) {
	var envelope dto.SqsEnvelope
	if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &envelope); err != nil {
		return nil, err
	}

	var pushMsg dto.PushMessage
	if err := json.Unmarshal([]byte(envelope.Message), &pushMsg); err != nil {
		return nil, err
	}

	return &pushMsg, nil
}
