package grpc

import (
	"context"
	"net"

	"push/common/lib"
	pb "push/linker/api/proto"
	"push/linker/internal/api/dto"
	"push/linker/internal/service"

	"go.uber.org/fx"
	"google.golang.org/grpc"
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
		Id:     uint(req.Id),
		Status: req.Status,
	}
	if err := s.service.UpdateMessageStatus(ctx, dto); err != nil {
		s.logger.Error(err)
	}

	return &pb.ResUpdateStatus{Reply: 1}, nil
}

// grpc.Server 생성
func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// gRPC 서버 시작 및 종료를 fx 라이프사이클 훅에 등록
func RegisterGRPCServer(lc fx.Lifecycle, grpcServer *grpc.Server, service pb.MessageServiceServer, log lib.Logger, env lib.Env) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":"+env.Linker.GrpcPort)
			if err != nil {
				return err
			}

			pb.RegisterMessageServiceServer(grpcServer, service)

			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					log.Errorf("gRPC server stopped: %v\n", err)
				}
			}()

			log.Debug("gRPC server started on :" + env.Linker.GrpcPort)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			log.Debug("gRPC server stopped gracefully")
			return nil
		},
	})
}

var Module = fx.Options(
	fx.Provide(NewGRPCServer),
	fx.Provide(NewMessageServiceServer),
	fx.Invoke(RegisterGRPCServer),
)
