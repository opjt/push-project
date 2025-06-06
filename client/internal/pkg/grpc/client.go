package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"push/common/lib"
	pb "push/dispatcher/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type sessionClient struct {
	client pb.SessionServiceClient
	logger lib.Logger
}

type SessionClient interface {
	Connect(context.Context, string, chan<- string) error
}

func NewSessionServiceClient(logger lib.Logger, lc fx.Lifecycle, env lib.Env) (SessionClient, error) {
	clientConn, err := grpc.NewClient("localhost:"+env.Dispatcher.SessionPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC: %w", err)
	}

	client := pb.NewSessionServiceClient(clientConn)

	c := &sessionClient{
		client: client,
		logger: logger,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing gRPC session client")
			return clientConn.Close()
		},
	})

	return c, nil
}

func (c *sessionClient) Connect(ctx context.Context, userID string, messageCh chan<- string) error {
	stream, err := c.client.Connect(ctx, &pb.ConnectRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("failed to connect to session stream: %w", err)
	}

	// 메시지 수신 고루틴
	go func() {
		defer close(messageCh)
		for {
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || status.Code(err) == codes.Canceled || status.Code(err) == codes.Unavailable {
					return
				}
				c.logger.Errorf("error receiving from stream: %v", err)
				return
			}
			if msg.GetMessage() == "__shutdown__" {
				return
			}

			// 받은 메시지를 채널로 전달
			messageCh <- msg.GetMessage()
		}
	}()

	return nil
}

var Module = fx.Options(
	fx.Provide(NewSessionServiceClient),
)
