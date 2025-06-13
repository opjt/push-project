package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"push/client/internal/tui/state"
	"push/common/lib"
	pb "push/dispatcher/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Message struct {
	MsgId uint64
	Title string
	Body  string
}

type sessionClient struct {
	client pb.SessionServiceClient
	logger lib.Logger
}

type SessionClient interface {
	Connect(context.Context, state.User, chan<- Message) error
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
			logger.Debug("closing gRPC session client")
			return clientConn.Close()
		},
	})

	return c, nil
}

func (c *sessionClient) Connect(ctx context.Context, user state.User, messageCh chan<- Message) error {
	stream, err := c.client.Connect(ctx, &pb.ConnectRequest{UserId: user.UserId, SessionId: user.SessionId})
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
			if msg.GetTitle() == "__shutdown__" {
				return
			}

			// 받은 메시지를 채널로 전달
			messageCh <- Message{
				MsgId: msg.GetMsgId(),
				Title: msg.GetTitle(),
				Body:  msg.GetBody(),
			}
		}
	}()

	return nil
}

var Module = fx.Options(
	fx.Provide(NewSessionServiceClient),
)
