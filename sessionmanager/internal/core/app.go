package core

import (
	"context"
	"push/common/lib/env"
	"push/common/lib/logger"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunServer(opt fx.Option) {
	log, _ := logger.NewLogger(env.NewEnv())

	app := fx.New(
		opt,
		fx.WithLogger(func() fxevent.Logger {
			return logger.NewFxLogger(log)
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Fatal("Failed to start app:", err)
	}

	// 시그널 대기
	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Error("Failed to stop app gracefully:", err)
	}
}
