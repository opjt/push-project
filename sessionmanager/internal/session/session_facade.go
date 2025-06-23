package session

import (
	"push/common/lib/logger"
	"push/linker/api/client"
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
func (r *SessionFacade) SendMessageToUser(pushDto *dto.Push) (int, error) {
	userId := pushDto.UserId
	sessionIDs := r.userSessionPool.GetSessionIDs(userId)
	sendCount := 0 // 전송된 메세지 수
	if len(sessionIDs) == 0 {
		//유저와 연결된 세션이 없을 경우.
		return sendCount, nil
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
			continue
		}
		sendCount++
	}
	return sendCount, nil
}
