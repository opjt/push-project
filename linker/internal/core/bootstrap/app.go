package bootstrap

import (
	"context"
	"push/common/lib"
	"push/linker/internal/api/router"
	"push/linker/internal/core"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunServer(opt fx.Option) {
	logger := lib.GetLogger()
	opts := fx.Options(
		fx.WithLogger(func() fxevent.Logger {
			return logger.GetFxLogger()
		}),
		fx.Invoke(Run()),
	)
	ctx := context.Background()
	app := fx.New(opt, opts)
	err := app.Start(ctx)
	defer app.Stop(ctx)
	if err != nil {
		logger.Fatal(err)
	}
}

func Run() any {
	return func(
		env lib.Env,
		engine core.Engine,
		route router.Routers,
		logger lib.Logger,

	) {
		route.Setup()

		logger.Info("Running server")

		if env.Linker.Port == "" {
			_ = engine.Gin.Run()
		} else {
			_ = engine.Gin.Run(":" + env.Linker.Port)
		}
	}
}
