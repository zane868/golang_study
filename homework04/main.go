package main

import (
	"fmt"
	"log"

	"github.com/zane868/golang_study/homework04/config"
	"github.com/zane868/golang_study/homework04/handler"
	"github.com/zane868/golang_study/homework04/middleware"
	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/service"
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
	postService := service.NewPostService(userService, db)
	userHandler := handler.NewUserHandler(userService, []byte(cfg.Jwt.Secret))
	postHandler := handler.NewPostHandler(postService)

	users_v1 := router.Group("/api/v1/users")
	{
		users_v1.POST("/register", userHandler.Register)
		users_v1.POST("/login", userHandler.Login)
	}

	post_v1 := router.Group("/api/v1/posts")
	post_v1.Use(middleware.Auth([]byte(cfg.Jwt.Secret)))
	{
		post_v1.POST("", postHandler.PublishBlog)
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
	if err := db.AutoMigrate(&model.User{}, &model.Post{}); err != nil {
		panic(fmt.Errorf("创建数据表失败：: %w", err))
	}

	return db
}
