package gin

import (
	"push/common/lib/logger"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	Gin      *gin.Engine
	ApiGroup *gin.RouterGroup
}

func NewEngine(logger *logger.GinLogger) Engine {
	gin.DefaultWriter = logger
	engine := gin.New()

	apiGroup := engine.Group("/api/v1")
	return Engine{Gin: engine, ApiGroup: apiGroup}
}
