package session

import (
	"fmt"
	"push/common/lib/env"
	"push/sessionmanager/api/client"

	"go.uber.org/fx"
)

type SessionClients map[string]client.SessionClient

func NewSessionClients(lc fx.Lifecycle, env env.Env) (SessionClients, error) {
	podRange := env.Pod.Index

	addresses := make([]string, podRange)
	for i := 0; i < podRange; i++ {
		addresses[i] = fmt.Sprintf("localhost:%d", env.Session.Port+i)

	}

	clients := make(SessionClients)

	for _, addr := range addresses {
		c, err := client.NewSessioneServiceClient(lc, addr)
		if err != nil {
			return nil, err
		}

		clients[addr] = c
	}

	return clients, nil
}
