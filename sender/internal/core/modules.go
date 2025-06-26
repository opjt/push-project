package core

import (
	"push/common"
	"push/linker/api/client"
	"push/sender/internal/service"
	"push/sender/internal/service/session"
	"push/sender/internal/sqs"
	"push/sessionmanager/sessionstore"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	common.CommonModules,
	fx.Provide(service.NewSenderService),
	fx.Provide(sqs.NewHandler),
	fx.Invoke(sqs.NewConsumer),
	client.Module,

	fx.Provide(session.NewSessionClients),

	fx.Provide(sessionstore.NewRedisClient),
	fx.Provide(sessionstore.NewReadRepository),
)
