package session

import (
	"fmt"
	pb "push/dispatcher/api/proto"
	"sync"
)

type InMemoryManager struct {
	sessions sync.Map
}

func NewInMemoryManager() Manager {
	return &InMemoryManager{
		sessions: sync.Map{},
	}
}

func (m *InMemoryManager) Add(userID string, stream pb.SessionService_ConnectServer) {
	m.sessions.Store(userID, &Session{
		UserID: userID,
		Stream: stream,
	})
}

func (m *InMemoryManager) Remove(userID string) {
	m.sessions.Delete(userID)
}

func (m *InMemoryManager) Get(userID string) (*Session, bool) {
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
func (m *InMemoryManager) SendTo(userID string, msg *pb.ServerMessage) error {
	val, ok := m.sessions.Load(userID)
	if !ok {
		return fmt.Errorf("no active session for user %s", userID)
	}
	session := val.(*Session)
	return session.Send(msg)
}

func (m *InMemoryManager) Len() int {
	length := 0
	m.sessions.Range(func(_, _ any) bool {
		length++
		return true
	})
	return length
}
