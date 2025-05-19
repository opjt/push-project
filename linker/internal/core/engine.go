package core

import (
	"push/common/lib"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	Gin      *gin.Engine
	ApiGroup *gin.RouterGroup
}

func NewEngine(logger lib.Logger) Engine {
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()

	apiGroup := engine.Group("/api/v1")
	return Engine{Gin: engine, ApiGroup: apiGroup}
}
