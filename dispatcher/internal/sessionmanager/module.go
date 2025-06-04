package sessionmanager

import (
	"push/dispatcher/internal/sessionmanager/grpc"

	"go.uber.org/fx"
)

var Module = fx.Options(

	grpc.Module,
)
