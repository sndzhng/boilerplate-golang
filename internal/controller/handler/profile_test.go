package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func beforeTestProfile(test *testing.T) (
	*usecasemock.MockAdmin,
	*usecasemock.MockUser,
	handler.Profile,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockAdminUsecase := usecasemock.NewMockAdmin(controller)
	mockUserUsecase := usecasemock.NewMockUser(controller)
	profileHandler := handler.NewProfileHandler(mockAdminUsecase, mockUserUsecase)

	return mockAdminUsecase, mockUserUsecase, profileHandler
}

func TestProfileGetAdminByToken(test *testing.T) {
	mockAdminUsecase, _, profileHandler := beforeTestProfile(test)

	path := "/admin/{context}/profile"
	id := uint64(1)
	admin := entity.Admin{
		ID: &id,
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
		username := "username"
		returnAdmin := entity.Admin{
			ID:       &id,
			Role:     &entity.Role{},
			RoleID:   &id,
			Username: &username,
		}

		mockAdminUsecase.EXPECT().Get(admin).Return(returnAdmin, nil)

		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, mockMiddlewareAuthorization, profileHandler.GetAdminByToken)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		encodedReturnAdmin, err := json.Marshal(returnAdmin)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedReturnAdmin), response.Body.String())
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminUsecase.EXPECT().Get(admin).Return(entity.Admin{}, util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, mockMiddlewareAuthorization, profileHandler.GetAdminByToken)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, profileHandler.GetAdminByToken)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})
}

func TestProfileGetUserByToken(test *testing.T) {
	_, mockUserUsecase, profileHandler := beforeTestProfile(test)

	path := "/{context}/profile"
	id := uint64(1)
	user := entity.User{ID: &id}

	mockMiddlewareAuthorization := func(ginContext *gin.Context) {
		claims := middleware.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				Subject: fmt.Sprint(id),
			},
		}
		ginContext.Set("claims", &claims)
	}

	test.Run("Success", func(test *testing.T) {
		username := "username"
		name := "name"
		phone := "0987654321"
		returnUser := entity.User{
			ID:              &id,
			Admin:           &entity.Admin{},
			AdminID:         &id,
			Username:        &username,
			Name:            &name,
			Phone:           &phone,
			IsResetPassword: new(bool),
		}

		mockUserUsecase.EXPECT().Get(user).Return(returnUser, nil)

		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, mockMiddlewareAuthorization, profileHandler.GetUserByToken)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		encodedUser, err := json.Marshal(returnUser)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedUser), response.Body.String())
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserUsecase.EXPECT().Get(user).Return(entity.User{}, util.Error{Code: http.StatusInternalServerError})

		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, mockMiddlewareAuthorization, profileHandler.GetUserByToken)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router := gin.Default()

		router.GET(path, profileHandler.GetUserByToken)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})
}
