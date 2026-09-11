package middleware

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zane868/golang_study/homework04/util"
)

// RequestLog 记录请求结果，不记录请求正文、查询参数或认证头。
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := rand.Text()
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()

		status := c.Writer.Status()
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		slog.Log(c.Request.Context(), level, "HTTP request completed",
			"request_id", requestID, "method", c.Request.Method,
			"route", route, "status", status,
			"duration_ms", float64(time.Since(start).Microseconds())/1000,
			"user_id", c.GetUint("userID"))
	}
}

// Recover 使用结构化日志记录 panic，不打印可能包含凭据的请求或 panic 内容。
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.ErrorContext(c.Request.Context(), "HTTP handler panic",
					"request_id", c.GetString("requestID"),
					"panic_type", fmt.Sprintf("%T", recovered),
					"stack", string(debug.Stack()))
				c.Abort()
				if !c.Writer.Written() {
					util.Error(c, http.StatusInternalServerError, "Internal server error")
				}
			}
		}()
		c.Next()
	}
}
