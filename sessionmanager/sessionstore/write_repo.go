package sessionstore

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	userPrefix    = "user_sessions"
	sessionPrefix = "session_location"
)

type WriteRepository interface {
	SaveSession(ctx context.Context, userID uint64, sessionID, podID string) error
	DeleteSession(ctx context.Context, userID uint64, sessionID string) error
}
type sessionRepository struct {
	rdb *redis.Client
}

func NewWriteRepository(rdb *redis.Client) WriteRepository {
	return &sessionRepository{rdb: rdb}
}

func (s *sessionRepository) SaveSession(ctx context.Context, userID uint64, sessionID, podID string) error {

	pipe := s.rdb.Pipeline()

	pipe.SAdd(ctx, fmt.Sprintf("%s:%d", userPrefix, userID), sessionID)
	pipe.Set(ctx, fmt.Sprintf("%s:%s", sessionPrefix, sessionID), podID, 0) // TTL 사용하지 않음

	_, err := pipe.Exec(ctx)
	return err
}

func (s *sessionRepository) DeleteSession(ctx context.Context, userID uint64, sessionID string) error {
	pipe := s.rdb.Pipeline()
	pipe.SRem(ctx, fmt.Sprintf("%s:%d", userPrefix, userID), sessionID)
	pipe.Del(ctx, fmt.Sprintf("%s:%s", sessionPrefix, sessionID))
	_, err := pipe.Exec(ctx)
	return err
}
