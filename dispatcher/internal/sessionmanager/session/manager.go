package session

import (
	"errors"
	"fmt"
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
type SessionManager interface {
	Add(userID string, stream pb.SessionService_ConnectServer)
	Remove(userID string)
	Get(userID string) (*Session, bool)
	SendTo(userID string, msg *pb.ServerMessage) error
	Len() int
}

type sessionManager struct {
	sessions sync.Map
}

func NewInMemoryManager() SessionManager {
	return &sessionManager{
		sessions: sync.Map{},
	}
}

func (m *sessionManager) Add(userID string, stream pb.SessionService_ConnectServer) {
	m.sessions.Store(userID, &Session{
		UserID: userID,
		Stream: stream,
	})
}

func (m *sessionManager) Remove(userID string) {
	m.sessions.Delete(userID)
}

func (m *sessionManager) Get(userID string) (*Session, bool) {
	value, ok := m.sessions.Load(userID)
	if !ok {
		return nil, false
	}
	sess, ok := value.(*Session)
	if !ok {
		return nil, false
	}
	return sess, true
}
func (m *sessionManager) SendTo(userID string, msg *pb.ServerMessage) error {
	val, ok := m.sessions.Load(userID)
	if !ok {
		return fmt.Errorf("no active session for user %s", userID)
	}
	session := val.(*Session)
	return session.Send(msg)
}

func (m *sessionManager) Len() int {
	length := 0
	m.sessions.Range(func(_, _ any) bool {
		length++
		return true
	})
	return length
}
