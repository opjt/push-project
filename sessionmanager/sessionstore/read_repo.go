package sessionstore

import (
	"context"
	"fmt"
	"push/sessionmanager/internal/dto"

	"github.com/redis/go-redis/v9"
)

type ReadRepository interface {
	GetUserSessions(ctx context.Context, userID uint64) ([]dto.SessionInfo, error)
}
type readRepository struct {
	rdb *redis.Client
}

func NewReadRepository(rdb *redis.Client) ReadRepository {
	return &readRepository{rdb: rdb}
}

func (s *readRepository) GetUserSessions(ctx context.Context, userID uint64) ([]dto.SessionInfo, error) {
	key := fmt.Sprintf("user_sessions:%d", userID)
	result, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]dto.SessionInfo, 0, len(result))
	for sessionID, podID := range result {
		sessions = append(sessions, dto.SessionInfo{
			SessionID: sessionID,
			PodID:     podID,
		})
	}

	return sessions, nil
}
