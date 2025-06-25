package sessionstore

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
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
	key := fmt.Sprintf("user_sessions:%d", userID)
	return s.rdb.HSet(ctx, key, sessionID, podID).Err()
}

func (s *sessionRepository) DeleteSession(ctx context.Context, userID uint64, sessionID string) error {
	key := fmt.Sprintf("user_sessions:%d", userID)
	return s.rdb.HDel(ctx, key, sessionID).Err()
}
