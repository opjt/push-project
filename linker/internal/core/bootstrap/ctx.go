package bootstrap

import "context"

type AppContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

// fx.App이 종료되면 cancel 되는 전역 context.
func NewAppContext() *AppContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &AppContext{
		Ctx:    ctx,
		Cancel: cancel,
	}
}
