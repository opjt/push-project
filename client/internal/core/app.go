package core

import (
	"context"
	"log"
	"push/client/internal/tui"
	"push/common/lib"

	tea "github.com/charmbracelet/bubbletea"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunApp() {

	// Bubble Tea 종료 알림 채널 생성
	teaDone := make(chan struct{})
	logger := lib.GetLogger()
	app := fx.New(
		Modules,
		fx.WithLogger(func() fxevent.Logger {
			return logger.GetFxLogger()
		}),
		fx.Provide(
			tui.NewLoginModel,
			tui.NewChatModel,
			tui.NewRootModel,
		),
		fx.Invoke(func(lc fx.Lifecycle, root *tui.RootModel) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						p := tea.NewProgram(root, tea.WithAltScreen())
						root.TeaProgram = p
						if _, err := p.Run(); err != nil {
							logger.Error("Bubble Tea exited:", err)
						}
						close(teaDone) // 종료 알림
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Debug("App is stopping...")
					return nil
				},
			})
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Bubble Tea 종료 대기
	<-teaDone

	// fx 앱 종료 트리거
	if err := app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
