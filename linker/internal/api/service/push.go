package service

import (
	"context"
	"push/common/lib"

	"push/linker/internal/api/dto"
	"push/linker/internal/repository"
)

type PushService interface {
	PostPush(context.Context, dto.CreateMessageDTO) (uint, error)
}
type pushService struct {
	logger         lib.Logger
	repository     repository.UserRepository
	messageService MessageService
}

func NewPushService(
	logger lib.Logger,
	repository repository.UserRepository,
	messageService MessageService,
) PushService {
	return &pushService{
		logger:         logger,
		repository:     repository,
		messageService: messageService,
	}
}

func (s *pushService) PostPush(ctx context.Context, dto dto.CreateMessageDTO) (uint, error) {

	msgId, err := s.messageService.createMessage(ctx, dto)
	if err != nil {
		return 0, err
	}

	return msgId, nil
}
