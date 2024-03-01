package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/usecase"
	"github.com/sndzhng/gin-template/internal/util"
)

type (
	Auth interface {
		AdminLogin(c *gin.Context)
		UserLogin(c *gin.Context)
		UserReset(c *gin.Context)
	}
	authHandler struct {
		authUsecase usecase.Auth
	}
)

func NewAuthHandler(authUsecase usecase.Auth) Auth {
	return &authHandler{authUsecase: authUsecase}
}

func (handler *authHandler) AdminLogin(ginContext *gin.Context) {
	login := entity.Login{}
	err := ginContext.ShouldBindJSON(&login)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	accessToken, err := handler.authUsecase.AdminLogin(login)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusOK, accessToken)
}

func (handler *authHandler) UserLogin(ginContext *gin.Context) {
	login := entity.Login{}
	err := ginContext.ShouldBindJSON(&login)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	accessToken, err := handler.authUsecase.UserLogin(login)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusOK, accessToken)
}

func (handler *authHandler) UserReset(ginContext *gin.Context) {
	reset := entity.Reset{}
	err := ginContext.ShouldBindJSON(&reset)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	subject, err := util.GetClaimSubject(ginContext)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}
	reset.ID = &subject

	err = handler.authUsecase.UserReset(reset)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.Status(http.StatusOK)
}
