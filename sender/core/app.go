package core

import (
	"context"
	"os"
	"os/signal"
	"push/common/lib"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunServer(opt fx.Option) {
	logger := lib.GetLogger()

	app := fx.New(
		opt,
		fx.WithLogger(func() fxevent.Logger {
			return logger.GetFxLogger()
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(ctx); err != nil {
		logger.Fatal("Failed to start app:", err)
	}

	// signal handling은 여기에 작성
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		app.Stop(context.Background())
	}()

	<-app.Done()
	logger.Info("app is done")
	if err := app.Stop(ctx); err != nil {
		logger.Fatal("Failed to stop app:", err)
	}
}
