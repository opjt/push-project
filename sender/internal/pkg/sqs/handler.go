package sqs

import (
	"context"
	"encoding/json"
	"push/common/lib"
	"push/sender/internal/dto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler interface {
	HandleMessage(ctx context.Context, msg types.Message) error
}

// handler_impl.go (혹은 handler.go 내부)
type handler struct {
	log lib.Logger
	// 다른 의존성 (DB, RPC 등) 주입 가능
}

func NewHandler(log lib.Logger) Handler {
	return &handler{log: log}
}

func (h *handler) HandleMessage(ctx context.Context, msg types.Message) error {
	pushMsg, err := parseSqsMessage(msg)
	if err != nil {
		h.log.Errorf("Failed to parse message: %v", err)
		return err
	}

	h.log.Infof("Received push message: %+v", pushMsg)
	// 이후 pushMsg 처리
	return nil
}

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
