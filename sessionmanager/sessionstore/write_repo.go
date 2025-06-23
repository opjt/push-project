package sessionstore

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type WriteRepository interface {
	SaveSession(ctx context.Context, userID, sessionID, podID string) error
}
type sessionRepository struct {
	rdb *redis.Client
}

func NewWriteRepository(rdb *redis.Client) WriteRepository {
	return &sessionRepository{rdb: rdb}
}

func (s *sessionRepository) SaveSession(ctx context.Context, userID, sessionID, podID string) error {

	pipe := s.rdb.Pipeline()

	pipe.SAdd(ctx, "user_sessions:"+userID, sessionID)
	pipe.Set(ctx, "session_location:"+sessionID, podID, 30*time.Second) // TTL 30초

	_, err := pipe.Exec(ctx)
	return err
}
