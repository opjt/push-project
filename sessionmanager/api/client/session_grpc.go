package client

import (
	"context"
	"fmt"
	pb "push/sessionmanager/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sessionClient struct {
	client pb.SessionServiceClient
}

type SessionClient interface {
	PushMessage(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error)
}

// grpc client 생성자
func NewSessioneServiceClient(lc fx.Lifecycle, podAddr string) (SessionClient, error) {

	clientConn, err := grpc.NewClient(podAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)

	}
	sessionServiceClient := pb.NewSessionServiceClient(clientConn)

	c := sessionClient{
		client: sessionServiceClient,
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

func (m *sessionClient) PushMessage(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	return m.client.PushMessage(ctx, req)
}

var Module = fx.Options(
	fx.Provide(NewSessioneServiceClient),
)
