package handler_test

import (
	"bytes"
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

func beforeTestAuth(test *testing.T) (
	*usecasemock.MockAuth,
	handler.Auth,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockAuthUsecase := usecasemock.NewMockAuth(controller)
	authHandler := handler.NewAuthHandler(mockAuthUsecase)

	return mockAuthUsecase, authHandler
}

func TestAdminLogin(test *testing.T) {
	mockAuthUsecase, authHandler := beforeTestAuth(test)

	path := "/admin/{context}/auth/login"
	username := "username"
	password := "password"
	login := entity.Login{
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		accessToken := entity.AccessToken{AccessToken: new(string)}

		mockAuthUsecase.EXPECT().AdminLogin(login).Return(accessToken, nil)

		body, err := json.Marshal(login)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, authHandler.AdminLogin)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		encodedAccessToken, err := json.Marshal(accessToken)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedAccessToken), response.Body.String())
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAuthUsecase.EXPECT().AdminLogin(login).Return(entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(login)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, authHandler.AdminLogin)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		body, err := json.Marshal(entity.Login{})
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, authHandler.AdminLogin)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}

func TestUserLogin(test *testing.T) {
	mockAuthUsecase, authHandler := beforeTestAuth(test)

	path := "/{context}/auth/login"
	username := "username"
	password := "password"
	login := entity.Login{
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		accessToken := entity.AccessToken{AccessToken: new(string)}
		mockAuthUsecase.EXPECT().UserLogin(login).Return(accessToken, nil)

		body, err := json.Marshal(login)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, authHandler.UserLogin)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)

		encodedAccessToken, err := json.Marshal(accessToken)
		assert.NoError(test, err)
		assert.Equal(test, string(encodedAccessToken), response.Body.String())
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAuthUsecase.EXPECT().UserLogin(login).Return(entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(login)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, authHandler.UserLogin)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		body, err := json.Marshal(entity.Login{})
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.POST(path, authHandler.UserLogin)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}

func TestUserReset(test *testing.T) {
	mockAuthUsecase, authHandler := beforeTestAuth(test)

	path := "/{context}/auth/reset"
	id := uint64(1)
	password := "password"
	reset := entity.Reset{
		ID:       &id,
		Password: &password,
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
		mockAuthUsecase.EXPECT().UserReset(reset).Return(nil)

		body, err := json.Marshal(reset)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PATCH(path, mockMiddlewareAuthorization, authHandler.UserReset)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusOK, response.Code)
	})

	test.Run("InternalError/UserReset", func(test *testing.T) {
		mockAuthUsecase.EXPECT().UserReset(reset).Return(util.Error{Code: http.StatusInternalServerError})

		body, err := json.Marshal(reset)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PATCH(path, mockMiddlewareAuthorization, authHandler.UserReset)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("InternalError/GetClaimSubject", func(test *testing.T) {
		body, err := json.Marshal(reset)
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PATCH(path, authHandler.UserReset)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusInternalServerError, response.Code)
	})

	test.Run("BadRequest", func(test *testing.T) {
		body, err := json.Marshal(entity.Reset{})
		assert.NoError(test, err)

		request := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(body))
		response := httptest.NewRecorder()
		router := gin.Default()

		router.PATCH(path, authHandler.UserReset)
		router.ServeHTTP(response, request)

		assert.Equal(test, http.StatusBadRequest, response.Code)
	})
}
