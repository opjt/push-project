package pkg

import (
	"push/linker/internal/pkg/awssns"
	"push/linker/internal/pkg/database"
	"push/linker/internal/pkg/grpc"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(database.NewDB),
	fx.Provide(awssns.NewPublisher),
	grpc.Module,
)
