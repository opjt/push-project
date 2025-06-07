package grpc

import (
	"fmt"
	"push/common/lib"
	pb "push/dispatcher/api/proto"
	"push/dispatcher/internal/sessionmanager/session"
	"time"
)

type sessionServiceServer struct {
	pb.UnimplementedSessionServiceServer
	sessions   session.Manager // SessionManager Interface
	logger     lib.Logger
	shutdownCh chan struct{}
}

func NewSessionServiceServer(logger lib.Logger, manager session.Manager) pb.SessionServiceServer {
	return &sessionServiceServer{
		sessions:   manager,
		logger:     logger,
		shutdownCh: make(chan struct{}),
	}
}

func (s *sessionServiceServer) Connect(req *pb.ConnectRequest, stream pb.SessionService_ConnectServer) error {
	userID := req.GetUserId()
	s.logger.Debugf("User connected: %s", userID)

	// 세션 추가
	s.sessions.Add(userID, stream)
	defer s.sessions.Remove(userID)

	// Connect 직후 첫 메시지 전송
	err := stream.Send(&pb.ServerMessage{
		Message: fmt.Sprintf("Welcome %s! [%s]", userID, time.Now().Format(time.RFC3339)),
	})
	if err != nil {
		s.logger.Errorf("Initial stream send error for %s: %v", userID, err)
		return err
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCh:
			// 종료 메시지 전송
			_ = stream.Send(&pb.ServerMessage{Message: "__shutdown__"})
			return nil

		case <-stream.Context().Done():
			s.logger.Debugf("User disconnected: %s", userID)
			return nil

		case <-ticker.C:
			err := s.sessions.SendTo(userID, &pb.ServerMessage{
				Message: fmt.Sprintf("Hello %s! [%s]", userID, time.Now().Format(time.RFC3339)),
			})
			if err != nil {
				s.logger.Errorf("Stream error for %s: %v", userID, err)
				return err
			}
		}
	}
}
