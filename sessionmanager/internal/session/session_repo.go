package session

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	rdb *redis.Client
}

func NewSessionRepository(rdb *redis.Client) *SessionRepository {
	return &SessionRepository{rdb: rdb}
}

func (s *SessionRepository) SaveSession(userID, sessionID, podID string) error {
	ctx := context.Background()
	pipe := s.rdb.Pipeline()

	pipe.SAdd(ctx, "user_sessions:"+userID, sessionID)
	pipe.Set(ctx, "session_location:"+sessionID, podID, 30) // TTL 30초

	_, err := pipe.Exec(ctx)
	return err
}
