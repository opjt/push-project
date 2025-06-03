package core

import (
	"context"
	"log"
	"push/client/internal/tui"

	tea "github.com/charmbracelet/bubbletea"

	"go.uber.org/fx"
)

func RunApp() {

	// Bubble Tea 종료 알림 채널 생성
	teaDone := make(chan struct{})

	app := fx.New(
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
							log.Println("Bubble Tea exited:", err)
						}
						close(teaDone) // 종료 알림
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Println("App is stopping... cleaning resources")
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
