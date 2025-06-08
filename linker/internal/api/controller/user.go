package controller

import (
	"net/http"
	"push/common/lib"
	"push/linker/dto"
	"push/linker/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	logger  lib.Logger
	service service.UserService
}

func NewUserController(logger lib.Logger, service service.UserService) UserController {
	return UserController{
		service: service,
		logger:  logger,
	}
}

func (u UserController) Login(c *gin.Context) {
	var req dto.AuthLoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warnf("invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

	ctx := c.Request.Context()

	loginResult, err := u.service.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	res := dto.AuthLoginRes(loginResult)

	c.JSON(http.StatusOK, res)
}
