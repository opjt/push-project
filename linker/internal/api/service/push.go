package service

import (
	"context"
	"fmt"
	"push/common/lib"

	"push/linker/internal/api/dto"
	"push/linker/internal/pkg/awssns"
	"push/linker/internal/repository"
)

type PushService interface {
	PostPush(context.Context, dto.CreateMessageDTO) (uint, error)
}
type pushService struct {
	logger         lib.Logger
	repository     repository.UserRepository
	messageService MessageService
	snsPublisher   awssns.Publisher
}

func NewPushService(
	logger lib.Logger,
	repository repository.UserRepository,
	messageService MessageService,
	snsPublisher awssns.Publisher,
) PushService {
	return &pushService{
		logger:         logger,
		repository:     repository,
		messageService: messageService,
		snsPublisher:   snsPublisher,
	}
}

func (s *pushService) PostPush(ctx context.Context, dto dto.CreateMessageDTO) (uint, error) {

	msgpubId, err := s.snsPublisher.Publish(ctx, dto.Content)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	s.logger.Debug("msgpubId:", msgpubId)
	msgId, err := s.messageService.createMessage(ctx, dto)
	if err != nil {
		return 0, err
	}

	return msgId, nil
}
