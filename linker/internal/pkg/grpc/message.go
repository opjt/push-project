package grpc

import (
	"context"
	"push/common/lib"
	pb "push/linker/api/proto"
	"push/linker/internal/api/dto"
	"push/linker/internal/service"
)

type messageServiceServer struct {
	pb.UnimplementedMessageServiceServer

	service service.MessageService
	logger  lib.Logger
}

func NewMessageServiceServer(service service.MessageService, logger lib.Logger) pb.MessageServiceServer {
	return &messageServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *messageServiceServer) UpdateStatus(ctx context.Context, req *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error) {

	dto := dto.UpdateMessageDTO{
		Id:       uint(req.Id),
		Status:   req.Status,
		SnsMsgId: req.SnsMsgId,
	}
	if err := s.service.UpdateMessageStatus(ctx, dto); err != nil {
		s.logger.Error(err)
	}

	return &pb.ResUpdateStatus{Reply: 1}, nil
}
