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

	createMessageDto := dto.CreateMessageDTO(msgdto)

	// DB에 저장 먼저 이후, mq 로직수행.
	msgId, err := s.messageService.createMessage(ctx, createMessageDto)
	if err != nil {
		return 0, err
	}

	snsBody := dto.SnsBody{
		MsgId:  msgId,
		Title:  msgdto.Title,
		Body:   msgdto.Content,
		UserId: msgdto.UserId,
	}

	//SNS 발행
	msgpubId, err := s.snsPublisher.Publish(ctx, snsBody)
	if err != nil {
		return 0, err
	}
	s.logger.Debug("msgpubId:", msgpubId)

	return msgId, nil
}
