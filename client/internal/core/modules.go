package core

import (
	"push/client/internal/pkg/grpc"
	"push/common"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	common.CommonModules,
	grpc.Module,
)
