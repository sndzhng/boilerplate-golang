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
	Admin interface {
		Create(c *gin.Context)
		DeleteByID(c *gin.Context)
		GetAll(c *gin.Context)
		GetByID(c *gin.Context)
		Initial(c *gin.Context)
		UpdateByID(c *gin.Context)
	}
	adminHandler struct {
		adminUsecase usecase.Admin
	}
)

func NewAdminHandler(adminUsecase usecase.Admin) Admin {
	return &adminHandler{adminUsecase: adminUsecase}
}

func (handler *adminHandler) Create(ginContext *gin.Context) {
	admin := entity.Admin{}
	err := ginContext.ShouldBindJSON(&admin)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	admin.PreventField()

	err = handler.adminUsecase.Create(admin)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusCreated, admin)
}

func (handler *adminHandler) DeleteByID(ginContext *gin.Context) {
	admin := entity.Admin{}
	id, err := strconv.ParseUint(ginContext.Param("id"), 10, 64)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	admin.ID = &id

	err = handler.adminUsecase.Delete(admin)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.Status(http.StatusOK)
}

func (handler *adminHandler) GetAll(ginContext *gin.Context) {
	adminFilter := entity.AdminFilter{}
	_ = ginContext.ShouldBindQuery(&adminFilter)

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

	admins, err := handler.adminUsecase.GetAll(&adminFilter, &sortOrder, &pagination)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(
		http.StatusOK,
		entity.AdminsWithNavigate{
			Admins:     admins,
			Pagination: pagination,
			SortOrder:  sortOrder,
		},
	)
}

func (handler *adminHandler) GetByID(ginContext *gin.Context) {
	id, err := strconv.ParseUint(ginContext.Param("id"), 10, 64)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	admin := entity.Admin{ID: &id}
	admin, err = handler.adminUsecase.Get(admin)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.JSON(http.StatusOK, admin)
}

func (handler *adminHandler) Initial(ginContext *gin.Context) {
	err := handler.adminUsecase.Initial()
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.Status(http.StatusCreated)
}

func (handler *adminHandler) UpdateByID(ginContext *gin.Context) {
	admin := entity.Admin{}
	_ = ginContext.ShouldBindJSON(&admin)
	admin.PreventField()

	id, err := strconv.ParseUint(ginContext.Param("id"), 10, 64)
	if err != nil {
		util.HandleError(ginContext, util.Error{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	admin.ID = &id

	err = handler.adminUsecase.Update(admin)
	if err != nil {
		util.HandleError(ginContext, err)
		return
	}

	ginContext.Status(http.StatusOK)
}
