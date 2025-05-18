package lib

import "go.uber.org/fx"

// Lib 모듈관리.
var Module = fx.Options(

	fx.Provide(NewEnv),
	fx.Provide(GetLogger),
)
