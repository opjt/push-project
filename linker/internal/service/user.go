package service

import (
	"context"
	"push/common/lib"
	"push/linker/dto"

	"github.com/google/uuid"

	"push/linker/internal/pkg/database"
	"push/linker/internal/repository"
)

type UserService interface {
	Login(ctx context.Context, user dto.AuthLoginReq) (dto.LoginResult, error)
}

// User Domain Service
type userService struct {
	logger     lib.Logger
	db         *database.MariaDB
	repository repository.UserRepository
}

func NewUserService(
	logger lib.Logger,
	db *database.MariaDB,
	repository repository.UserRepository,
) UserService {
	return &userService{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

func (s *userService) Login(ctx context.Context, loginReq dto.AuthLoginReq) (dto.LoginResult, error) {
	user, err := s.repository.GetUserByUsername(ctx, loginReq.Username)
	if err != nil {
		return dto.LoginResult{}, err
	}
	uuid := GenUuid()
	result := dto.LoginResult{
		Username:  user.Username,
		UserId:    user.ID,
		SessionId: uuid,
	}
	return result, nil
}

func GenUuid() string {
	id := uuid.New() // uuid v4
	return id.String()

}
