package sqs

import (
	"context"
	"encoding/json"
	"push/common/lib/logger"
	msgTypes "push/linker/types"
	"push/sender/internal/dto"
	"time"

	"push/linker/api/client"
	pb "push/linker/api/proto"
	sclient "push/sessionmanager/api/client"
	spb "push/sessionmanager/api/proto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler interface {
	HandleMessage(ctx context.Context, msg types.Message) error
}

type handler struct {
	log           *logger.Logger
	mclient       client.MessageClient
	sessionClient sclient.SessionClient
}

func NewHandler(log *logger.Logger, mclient client.MessageClient, sessionClient sclient.SessionClient) Handler {
	return &handler{
		log:           log,
		mclient:       mclient,
		sessionClient: sessionClient,
	}
}

func (h *handler) HandleMessage(ctx context.Context, msg types.Message) error {
	pushMsg, err := parseSqsMessage(msg)
	if err != nil {
		h.log.Errorf("Failed to parse message: %v", err)
		return err
	}

	h.log.Infof("Received push message: %+v", pushMsg)
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Linker에게 MessageStatus Update 요청
	if _, err = h.mclient.UpdateStatus(ctx, &pb.ReqUpdateStatus{Id: uint64(pushMsg.MsgID), Status: msgTypes.StatusSending, SnsMsgId: *msg.MessageId}); err != nil {
		h.log.Errorf("Failed to update message status: %v", err)
		return err
	}
	// SessionManager에게 메세지 전송 요청
	return h.sendPushMessage(pushMsg)
}

func (h *handler) sendPushMessage(pushMsg *dto.PushMessage) error {
	pushReq := spb.PushRequest{
		UserId:    uint64(pushMsg.UserID),
		SessionId: "",
		Message: &spb.ServerMessage{
			MsgId: uint64(pushMsg.MsgID),
			Title: pushMsg.Title,
			Body:  pushMsg.Body,
		},
	}
	result, err := h.sessionClient.PushMessage(context.Background(), &pushReq)
	if !result.Success {
		h.mclient.UpdateStatus(context.Background(), &pb.ReqUpdateStatus{Id: uint64(pushMsg.MsgID), Status: msgTypes.StatusDeferred}) // TODO: 에러처리 필요.
		return err
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
