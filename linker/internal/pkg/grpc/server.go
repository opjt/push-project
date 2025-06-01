package grpc

import (
	"context"
	"fmt"
	"net"

	"push/common/lib"
	pb "push/linker/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type messageServiceServer struct {
	pb.UnimplementedMessageServiceServer
}

func NewMessageServiceServer() pb.MessageServiceServer {
	return &messageServiceServer{}
}

func (s *messageServiceServer) UpdateStatus(ctx context.Context, req *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error) {
	// 여기서 원하는 비즈니스 로직 처리
	fmt.Printf("UpdateStatus called with id=%d, status=%s, sqsmsgId=%s\n", req.Id, req.Status, req.SqsmsgId)

	// 간단 응답 예시
	return &pb.ResUpdateStatus{Reply: "Status updated successfully"}, nil
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
					fmt.Printf("gRPC server stopped: %v\n", err)
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
