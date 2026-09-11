package util

import (
	"errors"
	"fmt"
	"log/slog"

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

func (e *AppError) Unwrap() error { return e.Err }

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		level := slog.LevelWarn
		if appErr.Code >= 500 {
			level = slog.LevelError
		}
		slog.Log(c.Request.Context(), level, "Application error",
			"request_id", c.GetString("requestID"), "code", appErr.Code, "error", err)
		Error(c, appErr.Code, appErr.Msg)
		return
	}
	// 未知错误，记录日志但不暴露给客户端
	slog.ErrorContext(c.Request.Context(), "Unhandled application error",
		"request_id", c.GetString("requestID"), "error", err)
	Error(c, 500, "Internal server error")
}
