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
	sessions SessionManager // sessionID -> stream

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
		logger:          logger,
		rpc:             rpc,
		writeRepository: writeRepository,
	}
}

// 세션 추가
func (r *SessionFacade) Add(userID uint64, sessionID string, stream pb.SessionService_ConnectServer) {
	r.sessions.Add(sessionID, stream)
	r.writeRepository.SaveSession(context.Background(), userID, sessionID, "localhost:50052") // TODO: podId 가져오기
}

// 세션 제거
func (r *SessionFacade) Remove(userID uint64, sessionID string) {
	r.sessions.Remove(sessionID)
	r.writeRepository.DeleteSession(context.Background(), userID, sessionID)
}

// 유저세션에 메시지 전송
func (r *SessionFacade) PushMessage(pushDto *dto.Push) error {
	sessionId := pushDto.SessionId
	stream, ok := r.sessions.Get(sessionId)
	if !ok { // 세션을 찾을 수 없을 경우
		r.writeRepository.DeleteSession(context.Background(), pushDto.UserId, sessionId) // TODO: 에러 처리
		return fmt.Errorf("session %s not found", sessionId)
	}

	err := stream.Send(&pb.ServerMessage{
		MsgId: pushDto.MsgId,
		Title: pushDto.Title,
		Body:  pushDto.Body,
	})

	if err != nil {
		return err
	}
	return nil
}
