package sessionmanager

import (
	"push/dispatcher/internal/sessionmanager/grpc"
	"push/dispatcher/internal/sessionmanager/session"

	"go.uber.org/fx"
)

var Module = fx.Options(
	session.Module,
	grpc.Module,
)
