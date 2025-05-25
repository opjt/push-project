package service

import (
	"context"
	"push/common/lib"
	"push/linker/internal/api/dto"
	"push/linker/internal/model"
	"push/linker/internal/pkg/database"
	"push/linker/internal/repository"
)

type MessageService interface {
	createMessage(context.Context, dto.CreateMessageDTO) (uint, error)
}

// 메세지를 생성하고 관리하는 서비스.
type messageService struct {
	logger     lib.Logger
	db         *database.MariaDB
	repository repository.MessageRepository
}

func NewMessageService(
	logger lib.Logger,
	db *database.MariaDB,
	repository repository.MessageRepository,
) MessageService {
	return &messageService{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

func (s *messageService) createMessage(ctx context.Context, dto dto.CreateMessageDTO) (uint, error) {
	msg := &model.Message{
		UserID:   dto.UserId,
		Title:    dto.Title,
		Content:  dto.Content,
		SnsMsgId: dto.SnsMsgId,
		Status:   model.STATUS_PENDING,
	}
	id, err := s.repository.CreateMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return id, nil
}
