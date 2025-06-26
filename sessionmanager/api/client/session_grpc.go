package client

import (
	"context"
	"fmt"
	"push/common/lib/env"
	"push/common/lib/logger"
	pb "push/sessionmanager/api/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sessionClient struct {
	client pb.SessionServiceClient
	logger *logger.Logger
}

type SessionClients map[string]SessionClient

type SessionClient interface {
	PushMessage(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error)
}

// grpc client 생성자
func NewSessioneServiceClient(logger *logger.Logger, lc fx.Lifecycle, env env.Env) (SessionClient, error) {

	clientConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", env.Session.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)

	}
	sessionServiceClient := pb.NewSessionServiceClient(clientConn)

	c := sessionClient{
		client: sessionServiceClient,
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

func NewSessionClients(logger *logger.Logger, lc fx.Lifecycle, env env.Env) (SessionClients, error) {
	podRange := env.Pod.Index

	addresses := make([]string, podRange)
	for i := 0; i < podRange; i++ {
		addresses[i] = fmt.Sprintf("localhost:%d", env.Session.Port+i)

	}

	clients := make(SessionClients)

	for _, addr := range addresses {
		clientConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to session service at %s: %w", addr, err)
		}

		sessionServiceClient := pb.NewSessionServiceClient(clientConn)
		logger.Info("Connected to session service at", addr)
		c := &sessionClient{
			client: sessionServiceClient,
			logger: logger,
		}

		// lifecycle에 클로저 등록
		conn := clientConn
		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return conn.Close()
			},
		})

		clients[addr] = c
	}

	return clients, nil
}

func (m *sessionClient) PushMessage(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	m.logger.Debugf("reqSessionClient  req: %v", req)
	return m.client.PushMessage(ctx, req)
}

var Module = fx.Options(
	fx.Provide(NewSessioneServiceClient),
)
