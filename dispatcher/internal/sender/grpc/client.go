package grpc

import (
	"context"
	"fmt"
	"push/common/lib"
	pb "push/linker/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO : 향후 공통서비스에서 해당 클라이언트를 제공하도록 개선 필요.
type messageClient struct {
	client pb.MessageServiceClient
	logger lib.Logger
}

type MessageClient interface {
	UpdateStatus(context.Context, *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error)
}

// grpc client 생성자
func NewMessageServiceClient(logger lib.Logger, lc fx.Lifecycle, env lib.Env) (MessageClient, error) {
	// Linker gRPC 연결.
	clientConn, err := grpc.NewClient("localhost:"+env.Linker.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)

	}
	messageServiceClient := pb.NewMessageServiceClient(clientConn)

	c := messageClient{
		client: messageServiceClient,
		logger: logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {

			return nil
		},
		OnStop: func(context.Context) error {
			clientConn.Close()
			return nil
		},
	})

	return &c, nil
}

func (m *messageClient) UpdateStatus(ctx context.Context, req *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error) {
	var result pb.ResUpdateStatus
	res, err := m.client.UpdateStatus(ctx, req)
	if err != nil {
		return &result, err
	}

	return res, nil
}

var Module = fx.Options(
	fx.Provide(NewMessageServiceClient),
)
