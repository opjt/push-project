package pkg

import (
	"push/linker/internal/pkg/database"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(database.NewDB),
)
