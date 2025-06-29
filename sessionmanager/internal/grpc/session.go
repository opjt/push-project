package grpc

import (
	"context"
	"fmt"
	"push/common/lib/logger"
	pb "push/sessionmanager/api/proto"
	"push/sessionmanager/internal/dto"
	"push/sessionmanager/internal/session"
	"time"
)

type sessionServiceServer struct {
	pb.UnimplementedSessionServiceServer
	manager    *session.SessionFacade
	logger     *logger.Logger
	shutdownCh chan struct{}
}

func NewSessionServiceServer(logger *logger.Logger, manager *session.SessionFacade) pb.SessionServiceServer {
	return &sessionServiceServer{
		logger:     logger,
		shutdownCh: make(chan struct{}),
		manager:    manager,
	}
}

// session에 메시지 전송.
func (s *sessionServiceServer) PushMessage(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	dto := &dto.Push{
		UserId:    req.GetUserId(),
		Title:     req.Message.GetTitle(),
		Body:      req.Message.GetBody(),
		MsgId:     req.Message.MsgId,
		SessionId: req.SessionId,
	}
	err := s.manager.PushMessage(dto)

	if err != nil {
		return &pb.PushResponse{Success: false}, err
	}

	return &pb.PushResponse{Success: true}, err

}

// 클라이언트에서 session 연결 메서드
func (s *sessionServiceServer) Connect(req *pb.ConnectRequest, stream pb.SessionService_ConnectServer) error {
	userId := req.GetUserId()
	sessionId := req.GetSessionId()
	s.logger.Debugf("User connected: %d - %s", userId, sessionId)

	// 세션 추가
	s.manager.Add(userId, sessionId, stream)

	// 연결 해제 시 정리
	defer func() {
		s.manager.Remove(userId, sessionId)
	}()

	// Connect 직후 첫 메시지 전송
	err := stream.Send(&pb.ServerMessage{
		Title: fmt.Sprintf("Welcome %s! [%s]", sessionId, time.Now().Format(time.RFC3339)),
		Body:  "session connected",
	})
	if err != nil {
		s.logger.Errorf("Initial stream send error for %s: %v", userId, err)
		return err
	}

	for {
		select {
		case <-s.shutdownCh:
			// 종료 메시지 전송
			_ = stream.Send(&pb.ServerMessage{Title: "__shutdown__"})
			return nil

		case <-stream.Context().Done():
			s.logger.Debugf("User disconnected: %s", sessionId)
			return nil
		}
	}
}
