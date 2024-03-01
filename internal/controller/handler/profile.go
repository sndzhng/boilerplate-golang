package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/usecase"
	"github.com/sndzhng/gin-template/internal/util"
)

type (
	Profile interface {
		GetAdminByToken(ginContext *gin.Context)
		GetUserByToken(ginContext *gin.Context)
	}
	profileHandler struct {
		adminUsecase usecase.Admin
		userUsecase  usecase.User
	}
)

func NewProfileHandler(
	adminUsecase usecase.Admin,
	userUsecase usecase.User,
) Profile {
	return &profileHandler{
		adminUsecase: adminUsecase,
		userUsecase:  userUsecase,
	}
}

func (handler *profileHandler) GetAdminByToken(ginContext *gin.Context) {
	subject, err := util.GetClaimSubject(ginContext)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	admin := entity.Admin{ID: &subject}
	admin, err = handler.adminUsecase.Get(admin)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusOK, admin)
}

func (handler *profileHandler) GetUserByToken(ginContext *gin.Context) {
	subject, err := util.GetClaimSubject(ginContext)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	user := entity.User{ID: &subject}
	user, err = handler.userUsecase.Get(user)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusOK, user)
}
