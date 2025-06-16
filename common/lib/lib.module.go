package lib

import (
	"push/common/lib/env"
	"push/common/lib/logger"

	"go.uber.org/fx"
)

// Lib 모듈관리.
var Module = fx.Options(

	fx.Provide(env.NewEnv),
	logger.Module,
)
