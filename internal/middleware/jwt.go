package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sndzhng/gin-template/internal/config"
	"github.com/sndzhng/gin-template/internal/entity"
)

type CustomClaims struct {
	jwt.StandardClaims
	Roles []entity.RoleName `json:"roles"`
}

func Authorization(ginContext *gin.Context) {
	tokenString := strings.TrimPrefix(ginContext.Request.Header.Get("Authorization"), "Bearer ")

	tokenJWT, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(tokenJWT *jwt.Token) (interface{}, error) {
			_, ok := tokenJWT.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("invalid token")
			}

			return []byte(config.JWT.Key), nil
		})
	if err != nil {
		ginContext.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, ok := tokenJWT.Claims.(*CustomClaims)
	if !ok {
		ginContext.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ginContext.Set("claims", tokenJWT.Claims)
}

func GenerateJWT(subject uint64, roles ...entity.RoleName) (string, error) {
	expireMinute, err := time.ParseDuration(os.Getenv("JWT_EXPIRE_MINUTE"))
	if err != nil {
		return "", err
	}

	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireMinute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   strconv.FormatUint(subject, 10),
		},
		Roles: roles,
	}
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenJWT.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyRoles(expectedRoles ...entity.RoleName) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		if ginContext.Keys["claims"] == nil {
			ginContext.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := ginContext.MustGet("claims").(*CustomClaims)
		for _, expectedRole := range expectedRoles {
			for _, userRole := range claims.Roles {
				if userRole == expectedRole {
					return
				}
			}
		}

		ginContext.AbortWithStatus(http.StatusUnauthorized)
	}
}
