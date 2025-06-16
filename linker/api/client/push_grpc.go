package client

import (
	"context"
	"fmt"
	"push/common/lib/env"
	"push/common/lib/logger"
	pb "push/linker/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type messageClient struct {
	client pb.MessageServiceClient
	logger *logger.Logger
}

type MessageClient interface {
	UpdateStatus(context.Context, *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error)
}

// grpc client 생성자
func NewMessageServiceClient(logger *logger.Logger, lc fx.Lifecycle, env env.Env) (MessageClient, error) {
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
