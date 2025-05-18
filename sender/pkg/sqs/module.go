package sqs

import "go.uber.org/fx"

var Module = fx.Options(

	fx.Provide(NewSQSClient),
	fx.Invoke(NewConsumer),
)
