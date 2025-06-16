package session

import (
	"context"
	"push/common/lib/logger"
	"push/linker/api/client"
	linkerpb "push/linker/api/proto"
	"push/linker/types"
	pb "push/sessionmanager/api/proto"
	"push/sessionmanager/internal/dto"
)

type SessionFacade struct {
	sessions        SessionManager  // sessionID -> stream
	userSessionPool UserSessionPool // userID -> []sessionID
	logger          *logger.Logger
	rpc             client.MessageClient
}

func NewSessionFacade(sessions SessionManager, userPool UserSessionPool, logger *logger.Logger, rpc client.MessageClient) *SessionFacade {
	return &SessionFacade{
		sessions:        sessions,
		userSessionPool: userPool,
		logger:          logger,
		rpc:             rpc,
	}
}

// 세션 추가
func (r *SessionFacade) Add(userID uint64, sessionID string, stream pb.SessionService_ConnectServer) {
	r.sessions.Add(sessionID, stream)
	r.userSessionPool.Add(userID, sessionID)
}

// 세션 제거
func (r *SessionFacade) Remove(userID uint64, sessionID string) {
	r.sessions.Remove(sessionID)
	r.userSessionPool.Remove(userID, sessionID)
}

// 유저에게 메시지 전송
func (r *SessionFacade) SendMessageToUser(pushDto *dto.Push) error {
	userId := pushDto.UserId
	sessionIDs := r.userSessionPool.GetSessionIDs(userId)
	if len(sessionIDs) == 0 {
		r.updateMessageStatus(context.Background(), pushDto.MsgId)
		return nil
	}

	for _, sid := range sessionIDs {
		stream, ok := r.sessions.Get(sid)
		if !ok {
			r.logger.Warnf("Session %s for user %s not found", sid, userId)
			continue
		}
		err := stream.Send(&pb.ServerMessage{MsgId: pushDto.MsgId, Title: pushDto.Title, Body: pushDto.Body})
		if err != nil {
			r.logger.Errorf("Failed to send message to session %s (user %s): %v", sid, userId, err)
		}
	}
	return nil
}

func (r *SessionFacade) updateMessageStatus(ctx context.Context, msgId uint64) {
	_, err := r.rpc.UpdateStatus(ctx, &linkerpb.ReqUpdateStatus{Id: msgId, Status: types.StatusDeferred})
	if err != nil {
		r.logger.Warnf("Failed to update message status: %v", err)
	}

}
