package sender

import (
	"push/sender/internal/sender/grpc"
	"push/sender/internal/sender/sqs"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(sqs.NewHandler),
	fx.Invoke(sqs.NewConsumer),
	grpc.Module,
)
