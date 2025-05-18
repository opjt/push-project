package router

import (
	"push/common/lib"
	"push/linker/api/controller"
	"push/linker/core"
)

type MessageRouter struct {
	logger            lib.Logger
	engine            core.Engine
	messageController controller.MessageController
}

func NewMessageRouter(
	logger lib.Logger,
	engine core.Engine,
	messageController controller.MessageController,

) Route {
	return MessageRouter{
		engine:            engine,
		logger:            logger,
		messageController: messageController,
	}
}

func (r MessageRouter) Setup() {
	msgRoutes := r.engine.ApiGroup.Group("/msg")
	{
		msgRoutes.POST("", r.messageController.Test)

	}

}
