package service

import (
	"context"
	"push/common/lib"

	"push/linker/internal/api/dto"
	"push/linker/internal/pkg/awssns"
	"push/linker/internal/repository"
)

type PushService interface {
	PostPush(context.Context, dto.PostPushDTO) (uint, error)
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

func (s *pushService) PostPush(ctx context.Context, msgdto dto.PostPushDTO) (uint, error) {

	snsBody := dto.SnsBody{
		Title:  msgdto.Title,
		Body:   msgdto.Content,
		UserId: msgdto.UserId,
	}

	msgpubId, err := s.snsPublisher.Publish(ctx, snsBody)
	if err != nil {
		return 0, err
	}
	s.logger.Debug("msgpubId:", msgpubId)

	createMessageDto := dto.CreateMessageDTO{
		Title:    msgdto.Title,
		Content:  msgdto.Content,
		UserId:   msgdto.UserId,
		SnsMsgId: msgpubId,
	}

	msgId, err := s.messageService.createMessage(ctx, createMessageDto)
	if err != nil {
		return 0, err
	}

	return msgId, nil
}
