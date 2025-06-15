package core

import (
	"push/common"
	"push/linker/api/client"
	"push/sender/internal/sqs"
	sclient "push/sessionmanager/api/client"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	common.CommonModules,
	fx.Provide(sqs.NewHandler),
	fx.Invoke(sqs.NewConsumer),
	client.Module,
	fx.Provide(sclient.NewSessioneServiceClient),
)
