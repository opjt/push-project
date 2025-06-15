package session

import (
	"sync"
)

type UserSessionPool interface {
	Add(userID uint64, sessionID string)
	Remove(userID uint64, sessionID string)
	GetSessionIDs(userID uint64) []string
}

type userSessionPool struct {
	userSessions sync.Map // map[userID]*sync.Map(sessionID -> struct{})
}

func NewInMemoryUserSessionStore() UserSessionPool {
	return &userSessionPool{}
}

func (s *userSessionPool) Add(userID uint64, sessionID string) {
	val, _ := s.userSessions.LoadOrStore(userID, &sync.Map{})
	sessionMap := val.(*sync.Map)
	sessionMap.Store(sessionID, struct{}{})
}

func (s *userSessionPool) Remove(userID uint64, sessionID string) {
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

func (s *userSessionPool) GetSessionIDs(userID uint64) []string {
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
