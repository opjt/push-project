package router

import (
	"push/common/lib"
	"push/linker/internal/api/controller"
	"push/linker/internal/pkg/gin"
)

type PushRouter struct {
	logger         lib.Logger
	engine         gin.Engine
	pushController controller.PushController
}

func NewPushRouter(
	logger lib.Logger,
	engine gin.Engine,
	pushController controller.PushController,

) Route {
	return PushRouter{
		engine:         engine,
		logger:         logger,
		pushController: pushController,
	}
}

func (r PushRouter) Setup() {
	pushRoutes := r.engine.ApiGroup.Group("/push")
	{
		pushRoutes.POST("/messages/users/:userid", r.pushController.PostPush)

	}

}
