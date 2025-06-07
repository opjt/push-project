package session

import (
	"errors"
	pb "push/dispatcher/api/proto"
	"sync"
)

type Session struct {
	UserID string
	Stream pb.SessionService_ConnectServer
	mu     sync.Mutex
}

func (s *Session) Send(msg *pb.ServerMessage) error {
	if s.Stream.Context().Err() != nil {
		return errors.New("stream closed")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Stream.Send(msg)
}

// SessionManager 인터페이스
type Manager interface {
	Add(userID string, stream pb.SessionService_ConnectServer)
	Remove(userID string)
	Get(userID string) (*Session, bool)
	SendTo(userID string, msg *pb.ServerMessage) error
	Len() int
}
