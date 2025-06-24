package session

import (
	"context"
	"fmt"
	"push/common/lib/logger"
	"push/linker/api/client"
	pb "push/sessionmanager/api/proto"
	"push/sessionmanager/internal/dto"
	"push/sessionmanager/sessionstore"
)

type SessionFacade struct {
	sessions        SessionManager  // sessionID -> stream
	userSessionPool UserSessionPool // userID -> []sessionID
	logger          *logger.Logger
	rpc             client.MessageClient
	writeRepository sessionstore.WriteRepository
}

func NewSessionFacade(
	sessions SessionManager, userPool UserSessionPool,
	logger *logger.Logger,
	rpc client.MessageClient,
	writeRepository sessionstore.WriteRepository,
) *SessionFacade {
	return &SessionFacade{
		sessions:        sessions,
		userSessionPool: userPool,
		logger:          logger,
		rpc:             rpc,
		writeRepository: writeRepository,
	}
}

// 세션 추가
func (r *SessionFacade) Add(userID uint64, sessionID string, stream pb.SessionService_ConnectServer) {
	r.sessions.Add(sessionID, stream)
	r.userSessionPool.Add(userID, sessionID)
	r.writeRepository.SaveSession(context.Background(), userID, sessionID, "pod1") // TODO: podId 가져오기
}

// 세션 제거
func (r *SessionFacade) Remove(userID uint64, sessionID string) {
	r.sessions.Remove(sessionID)
	r.userSessionPool.Remove(userID, sessionID)
	r.writeRepository.DeleteSession(context.Background(), userID, sessionID)
}

// 유저에게 메시지 전송
func (r *SessionFacade) SendMessageToUser(pushDto *dto.Push) (int, error) {
	userId := pushDto.UserId
	sessionIDs := r.userSessionPool.GetSessionIDs(userId)

	if len(sessionIDs) == 0 {
		r.logger.Debugf("No sessions found for user %d", userId)
		return 0, nil
	}

	var (
		sendCount int
		errs      []error
	)

	for _, sid := range sessionIDs {
		stream, ok := r.sessions.Get(sid)
		if !ok {
			errs = append(errs, fmt.Errorf("session %s not found", sid))
			continue
		}

		err := stream.Send(&pb.ServerMessage{
			MsgId: pushDto.MsgId,
			Title: pushDto.Title,
			Body:  pushDto.Body,
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to send to session %s: %w", sid, err))
			continue
		}
		r.logger.Debugf("Message sent to session %d - %s", userId, sid)
		sendCount++
	}

	// 여러 에러가 있을 경우 병합해서 리턴
	if len(errs) > 0 {
		return sendCount, fmt.Errorf("some messages failed to send: %v", errs)
	}

	return sendCount, nil
}
