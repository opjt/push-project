package session

import (
	"push/common/lib/env"
	"push/sessionmanager/api/client"

	"go.uber.org/fx"
)

type SessionClients map[string]client.SessionClient

func getPods(podCount int) []string {
	return []string{"localhost:50052"}
}
func NewSessionClients(lc fx.Lifecycle, env env.Env) (SessionClients, error) {
	podRange := env.Pod.Index

	addresses := getPods(podRange)

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
