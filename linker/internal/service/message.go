package service

import (
	"context"
	"push/common/lib/logger"
	"push/linker/internal/api/dto"
	"push/linker/internal/job/queue"
	"push/linker/types"

	"push/linker/internal/model"
	"push/linker/internal/repository"
	"time"
)

type MessageService interface {
	createMessage(context.Context, dto.CreateMessageDTO) (uint64, error)
	UpdateMessageStatus(context.Context, dto.UpdateMessageDTO) error
	UpdateMessagesStatus(context.Context, []uint64, dto.UpdateMessageField) error
	ReceiveMessage(context.Context, uint64) error
	UpdateStatusByJob(dto dto.UpdateMessageDTO) error
}

// 메세지를 생성하고 관리하는 서비스.
type messageService struct {
	logger       *logger.Logger
	repository   repository.MessageRepository
	queueManager *queue.JobQueueManager
}

func NewMessageService(
	logger *logger.Logger,
	repository repository.MessageRepository,
	queueManager *queue.JobQueueManager,

) MessageService {
	return &messageService{
		logger:       logger,
		repository:   repository,
		queueManager: queueManager,
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

func (s *messageService) UpdateMessagesStatus(ctx context.Context, ids []uint64, column dto.UpdateMessageField) error {
	if column.Status == types.StatusSent {
		now := time.Now()
		column.SentAt = &now
	}
	dto := dto.UpdateMessagesDTO{
		Ids:    ids,
		Column: column,
	}
	return s.repository.UpdateMessages(ctx, &dto)
}

func (s *messageService) ReceiveMessage(ctx context.Context, msgId uint64) error {

	return s.repository.ReceiveMessage(ctx, msgId)
}

func (s *messageService) UpdateStatusByJob(dto dto.UpdateMessageDTO) error {

	return s.queueManager.EnqUpdateStatus(dto)
}
