package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sndzhng/gin-template/internal/config"
	"github.com/sndzhng/gin-template/internal/controller/handler"
	"github.com/sndzhng/gin-template/internal/datastore"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/middleware"
	"github.com/sndzhng/gin-template/internal/repository"
	"github.com/sndzhng/gin-template/internal/usecase"
)

func SetupRouter() *gin.Engine {
	adminRepository := repository.NewAdminRepository(datastore.Postgresql)
	roleRepository := repository.NewRoleRepository(datastore.Postgresql)
	userRepository := repository.NewUserRepository(datastore.Postgresql)

	adminUsecase := usecase.NewAdminUsecase(adminRepository, roleRepository)
	authUsecase := usecase.NewAuthUsecase(adminRepository, userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)

	adminHandler := handler.NewAdminHandler(adminUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	profileHandler := handler.NewProfileHandler(adminUsecase, userUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	router := gin.Default()

	router.Use(
		cors.New(
			cors.Config{
				AllowCredentials: true,
				AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
				AllowOrigins:     []string{"*"},
				ExposeHeaders:    []string{"Content-Length"},
				MaxAge:           12 * time.Hour,
			},
		),
		func(ginContext *gin.Context) {
			ginContext.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			ginContext.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			ginContext.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			ginContext.Writer.Header().Set("Access-Control-Allow-Origin", "*")

			if ginContext.Request.Method == "OPTIONS" {
				ginContext.AbortWithStatus(http.StatusNoContent)
				return
			}

			ginContext.Next()
		},
	)

	router.NoRoute(func(ginContext *gin.Context) {
		ginContext.AbortWithStatus(http.StatusNotFound)
	})

	noAuthGroup := router.Group("")
	{
		noAuthGroup.POST(fmt.Sprintf("/admin/%s/admin/initial", config.Server.Context), adminHandler.Initial)
		noAuthGroup.POST(fmt.Sprintf("/admin/%s/auth/login", config.Server.Context), authHandler.AdminLogin)
		noAuthGroup.POST(fmt.Sprintf("/%s/auth/login", config.Server.Context), authHandler.UserLogin)
	}
	adminGroup := router.Group(fmt.Sprintf("/admin/%s", config.Server.Context), middleware.Authorization, middleware.VerifyRoles(entity.SuperAdminRoleName))
	{
		admin := adminGroup.Group("/admin")
		{
			admin.GET("", adminHandler.GetAll)
			admin.POST("", adminHandler.Create)
			admin.GET("/:id", adminHandler.GetByID)
			admin.PATCH("/:id", adminHandler.UpdateByID)
			admin.DELETE("/:id", adminHandler.DeleteByID)
		}
		profile := adminGroup.Group("/profile")
		{
			profile.GET("", profileHandler.GetAdminByToken)
		}
		user := adminGroup.Group("/user")
		{
			user.GET("", userHandler.GetAll)
			user.POST("", userHandler.Create)
			user.GET("/:id", userHandler.GetByID)
			user.PATCH("/:id", userHandler.UpdateByID)
			user.DELETE("/:id", userHandler.DeleteByID)
		}
	}
	userGroup := router.Group(fmt.Sprintf("/%s", config.Server.Context), middleware.Authorization, middleware.VerifyRoles(entity.UserRoleName))
	{
		auth := userGroup.Group("/auth")
		{
			auth.PATCH("/reset", authHandler.UserReset)
		}
		profile := userGroup.Group("/profile")
		{
			profile.GET("", profileHandler.GetUserByToken)
		}
	}

	return router
}
