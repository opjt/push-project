package pkg

import (
	"push/linker/internal/pkg/awssns"
	"push/linker/internal/pkg/awssqs"
	"push/linker/internal/pkg/database"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(database.NewDB),
	fx.Provide(awssns.NewPublisher),
	fx.Invoke(awssqs.NewConsumer),
)
