package router

import (
	"go.uber.org/fx"
)

type Route interface {
	Setup()
}

type Routers []Route

func NewRoutes(p struct {
	fx.In
	Routes []Route `group:"Routes"`
}) Routers {
	return Routers(p.Routes)
}

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewPushRouter,
			fx.ResultTags(`group:"Routes"`),
		),
	),
	fx.Provide(NewRoutes),
)

func (r Routers) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
