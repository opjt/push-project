package pkg

import (
	"push/sender/internal/pkg/sqs"

	"go.uber.org/fx"
)

var Module = fx.Options(

	fx.Provide(sqs.NewSQSClient),
	fx.Invoke(sqs.NewConsumer),
)
