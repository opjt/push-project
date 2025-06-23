package sessionstore

import "go.uber.org/fx"

var Module = fx.Options(

	fx.Provide(
		NewRedisClient,
		NewWriteRepository,
	),
)
