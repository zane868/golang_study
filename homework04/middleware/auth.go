package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zane868/golang_study/homework04/util"
)

func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// 验证 Token
		claims, err := util.ParseToken(authHeader, jwtSecret)
		if err != nil {
			util.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
