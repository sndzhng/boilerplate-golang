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
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/sndzhng/gin-template/internal/controller/handler"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/middleware"
	"github.com/sndzhng/gin-template/internal/util"
	usecasemock "github.com/sndzhng/gin-template/mock/usecase"
	"github.com/stretchr/testify/assert"
)

func beforeTestUser(test *testing.T) (
	*usecasemock.MockUser,
	handler.User,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockUserUsecase := usecasemock.NewMockUser(controller)
	userHandler := handler.NewUserHandler(mockUserUsecase)

	return mockUserUsecase, userHandler
}

func TestUserCreate(test *testing.T) {
	mockUserUsecase, userHandler := beforeTestUser(test)

	url := "/{context}/user"
	id := uint64(1)
	username := "username"
	password := "password"
	name := "name"
	phone := "0987654321"
	user := entity.User{
		AdminID:  &id,
		Username: &username,
		Password: &password,
		Name:     &name,
		Phone:    &phone,
	}

	mockMiddlewareAuthorization := func(ginContext *gin.Context) {
		claims := middleware.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				Subject: fmt.Sprint(id),
			},
		}
		ginContext.Set("claims", &claims)
	}

	test.Run("Success", func(test *testing.T) {
		mockUserUsecase.EXPECT().Create(user).Return(nil)

		body, err := json.Marshal(user)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(url, mockMiddlewareAuthorization, userHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusCreated, response.Code)
	})

	test.Run("InternalError/UserCreate", func(test *testing.T) {
		mockUserUsecase.EXPECT().Create(user).Return(util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(user)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(url, mockMiddlewareAuthorization, userHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("InternalError/GetClaimSubject", func(test *testing.T) {
		body, err := json.Marshal(user)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(url, userHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		body, err := json.Marshal(entity.User{})
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(url, userHandler.Create)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}

func TestUserDeleteByID(test *testing.T) {
	mockUserUsecase, userHandler := beforeTestUser(test)

	path := "/{context}/user/:id"
	id := uint64(1)
	user := entity.User{ID: &id}

	test.Run("Success", func(test *testing.T) {
		mockUserUsecase.EXPECT().Delete(user).Return(nil)

		request := httptest.NewRequest(http.MethodDelete, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE(path, userHandler.DeleteByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)
	})

	test.Run("InternalError", func(t *testing.T) {
		mockUserUsecase.EXPECT().Delete(user).Return(util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodDelete, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE(path, userHandler.DeleteByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.DELETE(path, userHandler.DeleteByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestUserGetAll(test *testing.T) {
	mockUserUsecase, userHandler := beforeTestUser(test)

	path := "/{context}/user"
	id := uint64(1)
	username := "username"
	name := "name"
	phone := "0987654321"
	userFilter := entity.UserFilter{
		User: entity.User{
			AdminID:  &id,
			Username: &username,
			Name:     &name,
			Phone:    &phone,
		},
	}
	sortOrder := entity.InitialSortOrder()
	pagination := entity.Pagination{
		Limit:  10,
		Offset: 0,
	}

	test.Run("Success", func(test *testing.T) {
		users := []entity.User{
			{
				ID:       &id,
				Admin:    &entity.Admin{},
				AdminID:  &id,
				Username: &username,
				Phone:    &phone,
			},
		}

		mockUserUsecase.EXPECT().GetAll(&userFilter, &sortOrder, &pagination).Return(users, nil)

		request := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf(
				"%s?admin_id=%d&username=%s&name=%s&phone=%s&limit=%d&offset=%d",
				path, id, username, name, phone, pagination.Limit, pagination.Offset,
			), nil,
		)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.GET(path, userHandler.GetAll)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		usersWithNavigate := entity.UsersWithNavigate{
			Users:      users,
			Pagination: pagination,
			SortOrder:  sortOrder,
		}
		encodedUsersWithNavigate, err := json.Marshal(usersWithNavigate)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedUsersWithNavigate), response.Body.String())
	})

	test.Run("InternalError", func(t *testing.T) {
		mockUserUsecase.EXPECT().GetAll(&userFilter, &sortOrder, &pagination).
			Return([]entity.User{}, util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf(
				"%s?admin_id=%d&username=%s&name=%s&phone=%s&limit=%d&offset=%d",
				path, id, username, name, phone, pagination.Limit, pagination.Offset,
			), nil,
		)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.GET(path, userHandler.GetAll)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.GET(path, userHandler.GetAll)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestUserGetByID(test *testing.T) {
	mockUserUsecase, userHandler := beforeTestUser(test)

	path := "/{context}/user/:id"
	id := uint64(1)
	user := entity.User{ID: &id}

	test.Run("Success", func(test *testing.T) {
		username := "username"
		password := "password"
		returnUser := entity.User{
			ID:       &id,
			Admin:    &entity.Admin{},
			AdminID:  &id,
			Username: &username,
			Password: &password,
		}

		mockUserUsecase.EXPECT().Get(user).Return(returnUser, nil)

		request := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.GET(path, userHandler.GetByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		encodedReturnUser, err := json.Marshal(returnUser)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedReturnUser), response.Body.String())
	})

	test.Run("InternalError", func(t *testing.T) {
		mockUserUsecase.EXPECT().Get(user).Return(entity.User{}, util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.GET(path, userHandler.GetByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.GET(path, userHandler.GetByID)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestUserUpdateByID(test *testing.T) {
	mockUserUsecase, userHandler := beforeTestUser(test)

	path := "/{context}/user/:id"
	id := uint64(1)
	username := "username"
	password := "password"
	user := entity.User{
		ID:       &id,
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		mockUserUsecase.EXPECT().Update(user).Return(nil)

		body, err := json.Marshal(user)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPut, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), bytes.NewReader(body))
		response := httptest.NewRecorder()

		router := gin.Default()
		router.PUT(path, userHandler.UpdateByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserUsecase.EXPECT().Update(user).Return(util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(user)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPut, strings.ReplaceAll(path, ":id", fmt.Sprintf("%d", id)), bytes.NewReader(body))
		response := httptest.NewRecorder()

		router := gin.Default()
		router.PUT(path, userHandler.UpdateByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		request := httptest.NewRequest(http.MethodPut, path, nil)
		response := httptest.NewRecorder()

		router := gin.Default()
		router.PUT(path, userHandler.UpdateByID)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}
