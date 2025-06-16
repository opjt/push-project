package core

import (
	"context"
	"net/http"
	"push/common/lib/env"
	"push/common/lib/logger"
	"push/linker/internal/api/router"
	"push/linker/internal/core/bootstrap"
	"push/linker/internal/pkg/gin"
	"time"

	"github.com/gin-contrib/pprof"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunServer(opt fx.Option) {
	log, _ := logger.NewLogger(env.NewEnv())
	opts := fx.Options(
		fx.WithLogger(func() fxevent.Logger {
			return logger.NewFxLogger(log)
		}),
		fx.Provide(bootstrap.NewAppContext),
		fx.Invoke(Run()),
	)
	app := fx.New(opt, opts)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// 시그널 대기
	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Error("Failed to stop app gracefully:", err)
	}
}

func Run() any {
	return func(
		lc fx.Lifecycle,
		env env.Env,
		engine gin.Engine,
		route router.Routers,
		logger *logger.Logger,
		appCtx *bootstrap.AppContext,
	) {
		route.Setup()

		logger.Info("Starting server on port: " + env.Linker.HttpPort)
		pprof.Register(engine.Gin)
		server := &http.Server{
			Addr:    ":" + env.Linker.HttpPort,
			Handler: engine.Gin,
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					// Run Gin server in goroutine
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						logger.Fatal("Failed to run server:", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Shutting down server gracefully...")
				appCtx.Cancel()
				return server.Shutdown(ctx)
			},
		})
	}
}
