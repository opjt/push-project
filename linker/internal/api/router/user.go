package router

import (
	"push/common/lib"
	"push/linker/internal/api/controller"
	"push/linker/internal/pkg/gin"
)

type UserRouter struct {
	logger         lib.Logger
	engine         gin.Engine
	userController controller.UserController
}

func NewUserRouter(
	logger lib.Logger,
	engine gin.Engine,
	userController controller.UserController,

) Route {
	return UserRouter{
		engine:         engine,
		logger:         logger,
		userController: userController,
	}
}

func (r UserRouter) Setup() {
	authRoute := r.engine.ApiGroup.Group("/auth")
	{
		authRoute.POST("/login", r.userController.Login)

	}

}
