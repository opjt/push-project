package service

import (
	"context"
	"push/linker/api/client"
	pb "push/linker/api/proto"
	msgtypes "push/linker/types"
	"push/sender/internal/dto"
	sclient "push/sessionmanager/api/client"
	spb "push/sessionmanager/api/proto"
)

type SenderService interface {
	PushMessage(context.Context, dto.PushMessage) error
	UpdateMessage(context.Context, dto.UpdateMessageStatus) error
}

type senderService struct {
	messageRpc client.MessageClient
	sessionRpc sclient.SessionClient
}

func NewSenderService(messageRpc client.MessageClient, sessionRpc sclient.SessionClient) SenderService {
	return &senderService{
		messageRpc: messageRpc,
		sessionRpc: sessionRpc,
	}
}

func (s *senderService) UpdateMessage(ctx context.Context, req dto.UpdateMessageStatus) error {
	_, err := s.messageRpc.UpdateStatus(
		ctx, &pb.ReqUpdateStatus{
			Id:       uint64(req.Id),
			Status:   req.Status,
			SnsMsgId: req.SnsMsgId,
		})
	return err

}

func (s *senderService) PushMessage(ctx context.Context, pushMsg dto.PushMessage) error {
	pushReq := spb.PushRequest{
		UserId:    uint64(pushMsg.UserID),
		SessionId: "",
		Message: &spb.ServerMessage{
			MsgId: uint64(pushMsg.MsgID),
			Title: pushMsg.Title,
			Body:  pushMsg.Body,
		},
	}
	result, err := s.sessionRpc.PushMessage(context.Background(), &pushReq)
	if err != nil {
		return err
	}
	if !result.Success {
		_, updateErr := s.messageRpc.UpdateStatus(context.Background(), &pb.ReqUpdateStatus{
			Id:     uint64(pushMsg.MsgID),
			Status: msgtypes.StatusDeferred,
		})
		return updateErr
	}
	return nil
}
