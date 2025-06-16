package grpc

import (
	"context"
	"net"

	"push/common/lib/env"
	"push/common/lib/logger"
	pb "push/linker/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// grpc.Server 생성
func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// gRPC 서버 시작 및 종료를 fx 라이프사이클 훅에 등록
func RegisterGRPCServer(lc fx.Lifecycle, grpcServer *grpc.Server, service pb.MessageServiceServer, log *logger.Logger, env env.Env) {
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
