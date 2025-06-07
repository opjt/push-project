package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewPushService),
	fx.Provide(NewUserService),
	fx.Provide(NewMessageService),
)
