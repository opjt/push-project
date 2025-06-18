package grpc

import (
	"context"
	"push/common/lib/logger"
	pb "push/linker/api/proto"
	"push/linker/internal/api/dto"
	"push/linker/internal/job/queue"
	"push/linker/internal/service"
)

type messageServiceServer struct {
	pb.UnimplementedMessageServiceServer

	service service.MessageService
	logger  *logger.Logger
}

func NewMessageServiceServer(service service.MessageService, queueManager *queue.JobQueueManager, logger *logger.Logger) pb.MessageServiceServer {
	return &messageServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *messageServiceServer) UpdateStatus(ctx context.Context, req *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error) {

	dto := dto.UpdateMessageDTO{
		Id:       req.Id,
		Status:   req.Status,
		SnsMsgId: req.SnsMsgId,
	}
	if err := s.service.UpdateStatusByJob(dto); err != nil {
		s.logger.Error(err)
	}

	return &pb.ResUpdateStatus{Reply: 1}, nil
}
