package service

import (
	"context"
	"push/common/lib"
	"push/linker/dto"
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
	user, err := s.repository.GetUserByUsername(ctx, loginReq.UserId)
	if err != nil {
		return dto.LoginResult{}, err
	}
	result := dto.LoginResult{
		UserId: user.Username,
	}
	return result, nil
}
