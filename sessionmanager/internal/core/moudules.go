package core

import (
	"push/common"
	"push/linker/api/client"
	"push/sessionmanager/internal/grpc"
	"push/sessionmanager/internal/session"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	common.CommonModules,
	grpc.Module,
	session.Module,
	fx.Provide(client.NewMessageServiceClient),
)
