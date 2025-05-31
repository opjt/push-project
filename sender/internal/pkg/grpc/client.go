package grpc

import (
	"context"
	"push/common/lib"
	pb "push/linker/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type messageClient struct {
	client pb.MessageServiceClient
	logger lib.Logger
}

type MessageClient interface {
	UpdateStatus(context.Context, *pb.ReqUpdateStatus) (*pb.ResUpdateStatus, error)
}

func NewMessageServiceClient(logger lib.Logger, lc fx.Lifecycle) MessageClient {

	clientConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("grpc client error") // TODO: 에러핸들링 수정 필요

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

	return &c
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
