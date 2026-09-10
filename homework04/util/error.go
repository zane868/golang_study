package util

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code int
	Msg  string
	Err  error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  msg,
	}
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.Code, appErr.Msg)
		return
	}
	// 未知错误，记录日志但不暴露给客户端
	Error(c, 500, "Internal server error")
}
