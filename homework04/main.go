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

	//实例化服务
	userService := service.NewUserService(db)
	postService := service.NewPostService(userService, db)
	userHandler := handler.NewUserHandler(userService, []byte(cfg.Jwt.Secret))
	postHandler := handler.NewPostHandler(postService)

	//注册路由
	router := gin.Default()

	//用户
	public_users_v1 := router.Group("/api/v1/users")
	{
		public_users_v1.POST("/register", userHandler.Register)
		public_users_v1.POST("/login", userHandler.Login)
	}

	//文章
	protected_post_v1 := router.Group("/api/v1/posts")
	protected_post_v1.Use(middleware.Auth([]byte(cfg.Jwt.Secret)))
	{
		protected_post_v1.GET("", postHandler.List)
		protected_post_v1.POST("", postHandler.PublishBlog)
		protected_post_v1.GET("/:id", postHandler.Get)
		protected_post_v1.PUT("/:id", postHandler.Update)
		protected_post_v1.DELETE("/:id", postHandler.Delete)
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
