package grpc

import (
	"context"
	"net"

	"push/common/lib"
	pb "push/dispatcher/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// grpc.Server 생성
func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

func RegisterGRPCServer(lc fx.Lifecycle, grpcServer *grpc.Server, service pb.SessionServiceServer, log lib.Logger, env lib.Env) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":"+env.Dispatcher.SessionPort)
			if err != nil {
				return err
			}

			pb.RegisterSessionServiceServer(grpcServer, service)

			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					log.Errorf("gRPC server stopped: %v\n", err)
				}
			}()

			log.Debug("gRPC server started on :" + env.Dispatcher.SessionPort)
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
	fx.Provide(NewSessionServiceServer),
	fx.Invoke(RegisterGRPCServer),
)
