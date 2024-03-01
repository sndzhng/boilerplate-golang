package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/usecase"
	"github.com/sndzhng/gin-template/internal/util"
)

type (
	User interface {
		Create(ginContext *gin.Context)
		DeleteByID(ginContext *gin.Context)
		GetAll(ginContext *gin.Context)
		GetByID(ginContext *gin.Context)
		UpdateByID(ginContext *gin.Context)
	}
	userHandler struct {
		userUsecase usecase.User
	}
)

func NewUserHandler(userUsecase usecase.User) User {
	return &userHandler{userUsecase: userUsecase}
}

func (handler *userHandler) Create(ginContext *gin.Context) {
	user := entity.User{}
	err := ginContext.ShouldBindJSON(&user)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	user.PreventField()

	subject, err := util.GetClaimSubject(ginContext)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	user.AdminID = &subject

	err = handler.userUsecase.Create(user)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusCreated, user)
}

func (handler *userHandler) DeleteByID(ginContext *gin.Context) {
	user := entity.User{}
	id, err := strconv.ParseUint(ginContext.Param("id"), 10, 64)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	user.ID = &id

	err = handler.userUsecase.Delete(user)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.Status(http.StatusOK)
}

func (handler *userHandler) GetAll(ginContext *gin.Context) {
	userFilter := entity.UserFilter{}
	_ = ginContext.ShouldBindQuery(&userFilter)

	sortOrder := entity.InitialSortOrder()
	err := ginContext.ShouldBindQuery(&sortOrder)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	if !sortOrder.Validate() {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: "invalid pagination sort order"})
		return
	}

	pagination := entity.Pagination{}
	err = ginContext.ShouldBindQuery(&pagination)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	users, err := handler.userUsecase.GetAll(&userFilter, &sortOrder, &pagination)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(
		http.StatusOK,
		entity.UsersWithNavigate{
			Users:      users,
			Pagination: pagination,
			SortOrder:  sortOrder,
		},
	)
}

func (handler *userHandler) GetByID(ginContext *gin.Context) {
	id, err := strconv.ParseUint(ginContext.Param("id"), 10, 64)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	user := entity.User{ID: &id}
	user, err = handler.userUsecase.Get(user)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusOK, user)
}

func (handler *userHandler) UpdateByID(ginContext *gin.Context) {
	user := entity.User{}
	_ = ginContext.ShouldBindJSON(&user)
	user.PreventField()

	id, err := strconv.ParseUint(ginContext.Param("id"), 10, 64)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	user.ID = &id

	err = handler.userUsecase.Update(user)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.Status(http.StatusOK)
}
