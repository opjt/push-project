package sqs

import (
	"context"
	"encoding/json"
	"push/common/lib"
	"push/dispatcher/internal/sender/dto"
	"push/dispatcher/internal/sender/grpc"
	"push/dispatcher/internal/sessionmanager/session"
	"time"

	pb "push/linker/api/proto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler interface {
	HandleMessage(ctx context.Context, msg types.Message) error
}

type handler struct {
	log           lib.Logger
	mclient       grpc.MessageClient
	sessionFacade *session.SessionFacade
}

func NewHandler(log lib.Logger, mclient grpc.MessageClient, sessionFacade *session.SessionFacade) Handler {
	return &handler{
		log:           log,
		mclient:       mclient,
		sessionFacade: sessionFacade,
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
	h.mclient.UpdateStatus(ctx, &pb.ReqUpdateStatus{Id: uint64(pushMsg.MsgID), Status: "sending"})

	return h.sendPushMessage(pushMsg)
}

func (h *handler) sendPushMessage(pushMsg *dto.PushMessage) error {

	return h.sessionFacade.SendMessageToUser(pushMsg)
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
