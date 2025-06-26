package service

import (
	"context"
	"push/linker/api/client"
	pb "push/linker/api/proto"
	msgtypes "push/linker/types"
	"push/sender/internal/dto"
	sclient "push/sessionmanager/api/client"
	spb "push/sessionmanager/api/proto"
	"push/sessionmanager/sessionstore"
)

type SenderService interface {
	PushMessage(context.Context, dto.PushMessage) error
	UpdateMessage(context.Context, dto.UpdateMessageStatus) error
}

type senderService struct {
	messageRpc     client.MessageClient
	sessionClients sclient.SessionClients
	sessionStore   sessionstore.ReadRepository
}

func NewSenderService(messageRpc client.MessageClient, sessionClients sclient.SessionClients, sessionStore sessionstore.ReadRepository) SenderService {
	return &senderService{
		messageRpc:     messageRpc,
		sessionClients: sessionClients,
		sessionStore:   sessionStore,
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

	if err := s.sendToUserSessions(ctx, &pushReq); err != nil {
		return err
	}

	return nil
}

func (s *senderService) sendToUserSessions(ctx context.Context, req *spb.PushRequest) error {

	sessions, err := s.sessionStore.GetUserSessions(ctx, req.UserId)

	if err != nil {
		return err
	}
	successCount := 0
	for _, session := range sessions {
		podAddr := session.PodID
		client := s.sessionClients[podAddr]
		req.SessionId = session.SessionID
		res, err := client.PushMessage(ctx, req)
		if err != nil {
			return err // TODO: error 처리 개선 필요.
		}
		if res.Success {
			successCount++
		}
	}

	if successCount == 0 {
		_, updateErr := s.messageRpc.UpdateStatus(context.Background(), &pb.ReqUpdateStatus{
			Id:     req.Message.MsgId,
			Status: msgtypes.StatusDeferred,
		})
		return updateErr
	}

	return nil
}
