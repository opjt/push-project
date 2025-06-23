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

func (s *sessionServiceServer) PushMessage(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	dto := &dto.Push{
		UserId: req.GetUserId(),
		Title:  req.Message.GetTitle(),
		Body:   req.Message.GetBody(),
		MsgId:  req.Message.MsgId,
	}
	sendCount, err := s.manager.SendMessageToUser(dto)

	resFlag := false
	if sendCount != 0 { // 만약 전송된 메시지 수가 0개가 아닐 경우 성공
		resFlag = true
	}

	return &pb.PushResponse{Success: resFlag}, err

}
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
