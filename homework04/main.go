package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zane868/golang_study/homework04/config"
	"github.com/zane868/golang_study/homework04/handler"
	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/service"
	"github.com/zane868/golang_study/homework04/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {

	//加载配置文件
	cfg := config.Load()

	// 初始化数据库
	db := initDb()

	//注册路由
	router := gin.Default()

	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService, []byte(cfg.Jwt.Secret))

	users_v1 := router.Group("/api/v1/users")
	{
		users_v1.POST("/register", userHandler.Register)
		users_v1.POST("/login", userHandler.Login)

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

func initDb() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("data/blogs.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	// 根据 User 结构自动创建或更新表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		panic(fmt.Errorf("创建数据表失败：: %w", err))
	}

	return db
}
