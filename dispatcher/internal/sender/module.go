package sender

import (
	"push/dispatcher/internal/sender/grpc"
	"push/dispatcher/internal/sender/sqs"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(sqs.NewHandler),
	fx.Invoke(sqs.NewConsumer),
	grpc.Module,
)
