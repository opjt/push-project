package core

import (
	"push/common"
	"push/linker/api/client"
	"push/sessionmanager/internal/grpc"
	"push/sessionmanager/internal/session"
	"push/sessionmanager/sessionstore"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	common.CommonModules,
	grpc.Module,
	sessionstore.Module,
	session.Module,
	fx.Provide(client.NewMessageServiceClient),
)
