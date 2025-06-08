package session

import (
	"sync"
)

type UserSessionPool interface {
	Add(userID, sessionID string)
	Remove(userID, sessionID string)
	GetSessionIDs(userID string) []string
}

type userSessionPool struct {
	userSessions sync.Map // map[userID]*sync.Map(sessionID -> struct{})
}

func NewInMemoryUserSessionStore() UserSessionPool {
	return &userSessionPool{}
}

func (s *userSessionPool) Add(userID, sessionID string) {
	val, _ := s.userSessions.LoadOrStore(userID, &sync.Map{})
	sessionMap := val.(*sync.Map)
	sessionMap.Store(sessionID, struct{}{})
}

func (s *userSessionPool) Remove(userID, sessionID string) {
	val, ok := s.userSessions.Load(userID)
	if !ok {
		return
	}
	sessionMap := val.(*sync.Map)
	sessionMap.Delete(sessionID)

	// userID 아래 세션이 비었으면 삭제
	empty := true
	sessionMap.Range(func(_, _ any) bool {
		empty = false
		return false
	})
	if empty {
		s.userSessions.Delete(userID)
	}
}

func (s *userSessionPool) GetSessionIDs(userID string) []string {
	val, ok := s.userSessions.Load(userID)
	if !ok {
		return nil
	}
	sessionMap := val.(*sync.Map)
	ids := make([]string, 0)
	sessionMap.Range(func(key, _ any) bool {
		ids = append(ids, key.(string))
		return true
	})
	return ids
}
