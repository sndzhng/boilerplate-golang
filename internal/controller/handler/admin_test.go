package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sndzhng/gin-template/internal/controller/handler"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/util"
	usecasemock "github.com/sndzhng/gin-template/mock/usecase"
	"github.com/stretchr/testify/assert"
)

func beforeTestAdmin(test *testing.T) (
	*usecasemock.MockAdmin,
	handler.Admin,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockAdminUsecase := usecasemock.NewMockAdmin(controller)
	adminHandler := handler.NewAdminHandler(mockAdminUsecase)

	return mockAdminUsecase, adminHandler
}

func TestAdminCreate(test *testing.T) {
	mockAdminUsecase, adminHandler := beforeTestAdmin(test)

	path := "/admin/{context}/admin"
	id := uint64(1)
	username := "username"
	password := "password"
	admin := entity.Admin{
		RoleID:   &id,
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		mockAdminUsecase.EXPECT().Create(admin).Return(nil)

		body, err := json.Marshal(admin)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, adminHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusCreated, response.Code)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminUsecase.EXPECT().Create(admin).Return(util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(admin)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, adminHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		body, err := json.Marshal(entity.Admin{})
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, adminHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}

func TestAdminDeleteByID(test *testing.T) {
	mockAdminUsecase, adminHandler := beforeTestAdmin(test)

	path := "/admin/{context}/admin/:id"
	id := uint64(1)
	admin := entity.Admin{ID: &id}

	test.Run("Success", func(test *testing.T) {
		mockAdminUsecase.EXPECT().Delete(admin).Return(nil)

		request := httptest.NewRequest(http.MethodDelete, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.DELETE(path, adminHandler.DeleteByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)
	})

	test.Run("InternalError", func(t *testing.T) {
		mockAdminUsecase.EXPECT().Delete(admin).Return(util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodDelete, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.DELETE(path, adminHandler.DeleteByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.DELETE(path, adminHandler.DeleteByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestAdminGetAll(test *testing.T) {
	mockAdminUsecase, adminHandler := beforeTestAdmin(test)

	path := "/admin/{context}/admin"
	id := uint64(1)
	username := "username"
	adminFilter := entity.AdminFilter{
		Admin: entity.Admin{
			RoleID:   &id,
			Username: &username,
		},
	}
	sortOrder := entity.InitialSortOrder()
	pagination := entity.Pagination{
		Limit:  10,
		Offset: 0,
	}

	test.Run("Success", func(test *testing.T) {
		admins := []entity.Admin{{ID: &id, Role: &entity.Role{}, RoleID: &id, Username: &username}}

		mockAdminUsecase.EXPECT().GetAll(&adminFilter, &sortOrder, &pagination).Return(admins, nil)

		request := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf(
				"%s?role_id=%d&name=%s&username=%s&limit=%d&offset=%d",
				path, id, username, username, pagination.Limit, pagination.Offset,
			), nil,
		)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, adminHandler.GetAll)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		adminsWithNavigate := entity.AdminsWithNavigate{
			Admins:     admins,
			Pagination: pagination,
			SortOrder:  sortOrder,
		}
		encodedAdminsWithNavigate, err := json.Marshal(adminsWithNavigate)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedAdminsWithNavigate), response.Body.String())
	})

	test.Run("InternalError", func(t *testing.T) {
		mockAdminUsecase.EXPECT().GetAll(&adminFilter, &sortOrder, &pagination).
			Return([]entity.Admin{}, util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf(
				"%s?role_id=%d&username=%s&limit=%d&offset=%d",
				path, id, username, pagination.Limit, pagination.Offset,
			), nil,
		)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, adminHandler.GetAll)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, adminHandler.GetAll)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestAdminGetByID(test *testing.T) {
	mockAdminUsecase, adminHandler := beforeTestAdmin(test)

	path := "/admin/{context}/admin/:id"
	id := uint64(1)
	admin := entity.Admin{ID: &id}

	test.Run("Success", func(test *testing.T) {
		username := "username"
		password := "password"
		returnAdmin := entity.Admin{
			ID:       &id,
			Role:     &entity.Role{},
			RoleID:   &id,
			Username: &username,
			Password: &password,
		}

		mockAdminUsecase.EXPECT().Get(admin).Return(returnAdmin, nil)

		request := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, adminHandler.GetByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		encodedReturnAdmin, err := json.Marshal(returnAdmin)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedReturnAdmin), response.Body.String())
	})

	test.Run("InternalError", func(t *testing.T) {
		mockAdminUsecase.EXPECT().Get(admin).Return(entity.Admin{}, util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, adminHandler.GetByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, adminHandler.GetByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestAdminInitial(test *testing.T) {
	mockAdminUsecase, adminHandler := beforeTestAdmin(test)

	path := "/admin/{context}/admin/initial"

	test.Run("Success", func(t *testing.T) {
		mockAdminUsecase.EXPECT().Initial().Return(nil)

		request := httptest.NewRequest(http.MethodPost, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, adminHandler.Initial)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	test.Run("InternalError", func(t *testing.T) {
		mockAdminUsecase.EXPECT().Initial().Return(util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodPost, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, adminHandler.Initial)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestAdminUpdateByID(test *testing.T) {
	mockAdminUsecase, adminHandler := beforeTestAdmin(test)

	path := "/admin/{context}/admin/:id"
	id := uint64(1)
	username := "username"
	password := "password"
	admin := entity.Admin{
		ID:       &id,
		RoleID:   &id,
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		mockAdminUsecase.EXPECT().Update(admin).Return(nil)

		body, err := json.Marshal(admin)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPut, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PUT(path, adminHandler.UpdateByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminUsecase.EXPECT().Update(admin).Return(util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(admin)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPut, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PUT(path, adminHandler.UpdateByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		request := httptest.NewRequest(http.MethodPut, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PUT(path, adminHandler.UpdateByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}
