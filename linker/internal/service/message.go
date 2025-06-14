package service

import (
	"context"
	"push/common/lib"
	"push/linker/internal/api/dto"
	"push/linker/internal/model"
	"push/linker/internal/pkg/database"
	"push/linker/internal/repository"
	"push/linker/types"
	"time"
)

type MessageService interface {
	createMessage(context.Context, dto.CreateMessageDTO) (uint64, error)
	UpdateMessageStatus(context.Context, dto.UpdateMessageDTO) error
	ReceiveMessage(context.Context, uint64) error
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

func (s *messageService) createMessage(ctx context.Context, dto dto.CreateMessageDTO) (uint64, error) {
	msg := &model.Message{
		UserID:  dto.UserId,
		Title:   dto.Title,
		Content: dto.Content,
		Status:  types.StatusPending,
	}
	id, err := s.repository.CreateMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *messageService) UpdateMessageStatus(ctx context.Context, dto dto.UpdateMessageDTO) error {

	msg := &model.Message{
		ID:       dto.Id,
		Status:   dto.Status,
		SnsMsgId: dto.SnsMsgId,
	}
	if dto.Status == types.StatusSent {
		now := time.Now()
		msg.SentAt = &now
	}
	return s.repository.UpdateMessage(ctx, msg)
}

func (s *messageService) ReceiveMessage(ctx context.Context, msgId uint64) error {

	return s.repository.ReceiveMessage(ctx, msgId)
}
