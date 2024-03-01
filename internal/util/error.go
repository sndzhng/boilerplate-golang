package util

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sndzhng/gin-template/internal/common"
	"github.com/sndzhng/gin-template/internal/config"
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func HandleError(ginContext *gin.Context, err error) {
	switch e := err.(type) {
	case Error:
		if e.Message != "" {
			log.Println(e.Message)
		}

		if e.Message != "" && config.Environment != common.Environment.Production {
			ginContext.AbortWithStatusJSON(e.Code, e.Message)
		} else {
			ginContext.AbortWithStatus(e.Code)
		}

		return
	default:
		log.Println(err.Error())
		ginContext.AbortWithError(500, err)

		return
	}
}

func NewError(code int, message string) Error {
	return Error{Code: code, Message: message}
}
