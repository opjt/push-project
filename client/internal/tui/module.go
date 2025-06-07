package tui

import (
	"push/client/internal/tui/state"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		state.NewUser,
		NewLoginModel,
		NewChatModel,
		NewRootModel,
	),
)
