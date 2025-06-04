package grpc

import (
	"fmt"
	"log"
	"push/common/lib"
	pb "push/dispatcher/api/proto"
	"sync"
	"time"
)

type sessionServiceServer struct {
	pb.UnimplementedSessionServiceServer
	sessions sync.Map
	logger   lib.Logger
}
type Session struct {
	UserID string
	Stream pb.SessionService_ConnectServer
}

func NewSessionServiceServer(logger lib.Logger) pb.SessionServiceServer {
	return &sessionServiceServer{
		sessions: sync.Map{},
		logger:   logger,
	}
}

func (s *sessionServiceServer) Connect(req *pb.ConnectRequest, stream pb.SessionService_ConnectServer) error {
	userID := req.GetUserId()

	log.Printf("User connected: %s", userID)

	s.sessions.Store(userID, &Session{
		UserID: userID,
		Stream: stream,
	})

	defer s.sessions.Delete(userID)

	for {
		// ping 메시지 전송 (stream 테스트용)
		err := stream.Send(&pb.ServerMessage{
			Message: fmt.Sprintf("Hello %s! [%s]", userID, time.Now().Format(time.RFC3339)),
		})
		if err != nil {
			log.Printf("Stream error for %s: %v", userID, err)
			return err
		}

		time.Sleep(5 * time.Second)
	}
}
