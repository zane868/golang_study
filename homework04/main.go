package main

import (
	"fmt"
	"homework01/homework04/config"
	"homework01/homework04/model"
	"homework01/homework04/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	//加载配置文件
	cfg := config.Load()

	token, _ := util.GenerateToken([]byte(cfg.Jwt.Secret), 556, "hello")

	fmt.Println(token)

	c, _ := util.ParseToken(token, []byte(cfg.Jwt.Secret))

	fmt.Println(c)

	//注册路由
	router := gin.Default()

	users_v1 := router.Group("/api/v1/users")
	{
		users_v1.POST("/register", func(ctx *gin.Context) {
			var u model.User
			if err := ctx.ShouldBindJSON(&u); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if u.UserName != "manu" || u.Password != "123" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		})

		users_v1.POST("/login", func(ctx *gin.Context) {
			var u model.User
			if err := ctx.ShouldBindJSON(&u); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			token, _ := util.GenerateToken([]byte(cfg.Jwt.Secret), 5, u.UserName)
			util.Success(ctx, gin.H{
				"token": token,
			})
		})

	}

	post_v1 := router.Group("/api/v1/post")
	{
		post_v1.POST("", func(ctx *gin.Context) {
			authHeader := ctx.GetHeader("Authorization")
			// 验证 Token
			claims, err := util.ParseToken(authHeader, []byte(cfg.Jwt.Secret))
			if err != nil {
				util.Error(ctx, http.StatusUnauthorized, "Invalid token")
				ctx.Abort()
				return
			}
			fmt.Println(claims)
		})
	}

	//启动服务
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	router.Run(addr) // listens on 0.0.0.0:8080 by default
}
