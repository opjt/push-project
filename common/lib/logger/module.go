package logger

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewLogger),
	fx.Provide(NewGormLogger),
	fx.Provide(NewGinLogger),
)
