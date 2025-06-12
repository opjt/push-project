package session

import (
	"push/common/lib"
	pb "push/dispatcher/api/proto"
)

type SessionFacade struct {
	sessions        SessionManager  // sessionID -> stream
	userSessionPool UserSessionPool // userID -> []sessionID
	logger          lib.Logger
}

func NewSessionFacade(sessions SessionManager, userPool UserSessionPool, logger lib.Logger) *SessionFacade {
	return &SessionFacade{
		sessions:        sessions,
		userSessionPool: userPool,
		logger:          logger,
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
func (r *SessionFacade) SendMessageToUser(userID uint64, title, body string) error {
	sessionIDs := r.userSessionPool.GetSessionIDs(userID)
	if len(sessionIDs) == 0 {
		r.logger.Infof("No active sessions for user %s", userID)
		return nil
	}

	for _, sid := range sessionIDs {
		stream, ok := r.sessions.Get(sid)
		if !ok {
			r.logger.Warnf("Session %s for user %s not found", sid, userID)
			continue
		}
		err := stream.Send(&pb.ServerMessage{Title: title, Body: body})
		if err != nil {
			r.logger.Errorf("Failed to send message to session %s (user %s): %v", sid, userID, err)
		}
	}
	return nil
}
