package service

import (
	"push/common/lib"
	"push/linker/internal/pkg/database"
	"push/linker/internal/repository"

	"github.com/gin-gonic/gin"
)

type PushService struct {
	logger     lib.Logger
	db         *database.MariaDB
	repository repository.UserRepository
}

func NewPushService(
	logger lib.Logger,
	db *database.MariaDB,
	repository repository.UserRepository,
) PushService {
	return PushService{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

func (s PushService) Test(c *gin.Context) error {
	if err := s.db.Ping(); err != nil {
		s.logger.Error("DB ping failed: %v", err)
		return err

	}
	ctx := c.Request.Context()
	usr, err := s.repository.GetUserByID(ctx, 1)
	if err != nil {
		s.logger.Error("GetUserByID failed: %v", err)
	}
	s.logger.Info("User: %v", usr)
	return nil
}
