package controller

import (
	"net/http"
	commondto "push/common/dto"
	"push/common/lib/logger"
	"push/linker/dto"
	"push/linker/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	logger  *logger.Logger
	service service.UserService
}

func NewUserController(logger *logger.Logger, service service.UserService) UserController {
	return UserController{
		service: service,
		logger:  logger,
	}
}

// 클라이언트에서 사용자 로그인을 위해 사용하는 endpoint
// POST /api/v1/auth/login
// 추후 jwt 를 통해 클라이언트 -> linker로 메세지 변경 요청시 인증/인가 시스템 구현 필요.
func (u UserController) Login(c *gin.Context) {
	var req dto.AuthLoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warnf("invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

	ctx := c.Request.Context()
	res := commondto.CommonResponse[dto.AuthLoginRes]{}
	loginResult, err := u.service.Login(ctx, req)
	if err != nil {
		res.Error = err.Error()
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	res.Data = dto.AuthLoginRes(loginResult)
	c.JSON(http.StatusOK, res)
}
