package router

import (
	"push/common/lib"
	"push/linker/internal/api/controller"
	"push/linker/internal/pkg/gin"
)

type MessageRouter struct {
	logger            lib.Logger
	engine            gin.Engine
	messageController controller.MessageController
}

func NewMessageRouter(
	logger lib.Logger,
	engine gin.Engine,
	messageController controller.MessageController,

) Route {
	return MessageRouter{
		engine:            engine,
		logger:            logger,
		messageController: messageController,
	}
}

func (r MessageRouter) Setup() {
	msgRoutes := r.engine.ApiGroup.Group("/messages")
	{
		msgRoutes.PATCH("/:msgid/status", r.messageController.UpdateStatus)

	}

}
