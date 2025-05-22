package service

import (
	"context"
	"fmt"
	"push/common/lib"
	"push/linker/internal/api/dto"
	"push/linker/internal/model"
	"push/linker/internal/pkg/database"
	"push/linker/internal/repository"
)

type MessageService struct {
	logger     lib.Logger
	db         *database.MariaDB
	repository repository.MessageRepository
}

func NewMessageService(
	logger lib.Logger,
	db *database.MariaDB,
	repository repository.MessageRepository,
) MessageService {
	return MessageService{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

func (s MessageService) Test(ctx context.Context) error {

	createDto := dto.CreateMessageDTO{
		UserID:  1,
		Content: "hi",
	}
	res, err := s.createMessage(ctx, createDto)
	if err != nil {
		return nil
	}
	fmt.Println(res)
	return nil
}

func (s MessageService) createMessage(ctx context.Context, dto dto.CreateMessageDTO) (uint, error) {
	msg := &model.Message{
		UserID:  dto.UserID,
		Content: dto.Content,
		Status:  "pending", // 기본 상태
	}
	id, err := s.repository.CreateMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return id, nil
}
