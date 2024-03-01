package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/middleware"
)

func GetClaimSubject(ginContext *gin.Context) (uint64, error) {
	if ginContext.Keys["claims"] == nil {
		return 0, errors.New("claims not found")
	}

	claims := ginContext.MustGet("claims").(*middleware.CustomClaims)
	subject, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return subject, nil
}

func GetClaimRoles(ginContext *gin.Context) ([]entity.RoleName, error) {
	if ginContext.Keys["claims"] == nil {
		return []entity.RoleName{}, errors.New("claims not found")
	}

	claims := ginContext.MustGet("claims").(*middleware.CustomClaims)
	roles := claims.Roles

	return roles, nil
}

func ModifyRequestBody(ginContext *gin.Context, modifyMap map[string]interface{}) error {
	bodyBytes, err := io.ReadAll(ginContext.Request.Body)
	if err != nil {
		return err
	}

	bodyMap := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &bodyMap)
	if err != nil {
		return err
	}

	for key, value := range modifyMap {
		bodyMap[key] = value
	}
	bodyBytes, err = json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	ginContext.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

func ModifyRequestParams(ginContext *gin.Context, modifyMap map[string]interface{}) {
	queryParams := ginContext.Request.URL.Query()

	for key, value := range modifyMap {
		queryParams.Set(key, fmt.Sprint(value))
	}

	ginContext.Request.URL.RawQuery = queryParams.Encode()
}
