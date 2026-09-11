package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zane868/golang_study/homework04/config"
	"github.com/zane868/golang_study/homework04/handler"
	"github.com/zane868/golang_study/homework04/middleware"
	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ServiceContext struct {
	UserService    *service.UserService
	PostService    *service.PostService
	CommentService *service.CommentService
	UserHandler    *handler.UserHandler
	PostHandler    *handler.PostHandler
	CommentHandler *handler.CommentHandler
}

//go:embed web/index.html
var indexHTML []byte

func main() {

	//加载配置文件
	cfg := config.Load()

	//初始化数据库
	db := initDb()

	//实例化服务
	serviceContext := NewContext(db, cfg)

	//注册路由
	router := regRouter(serviceContext, cfg)

	//启动服务
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	router.Run(addr) // listens on 0.0.0.0:8080 by default
}

func NewContext(db *gorm.DB, cfg *config.Config) *ServiceContext {
	userService := service.NewUserService(db)
	postService := service.NewPostService(userService, db)
	commentService := service.NewCommentService(db)

	return &ServiceContext{
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,

		UserHandler:    handler.NewUserHandler(userService, []byte(cfg.Jwt.Secret)),
		PostHandler:    handler.NewPostHandler(postService),
		CommentHandler: handler.NewCommentHandler(commentService),
	}
}

func regRouter(sc *ServiceContext, cfg *config.Config) *gin.Engine {
	//注册路由
	router := gin.Default()
	router.GET("/index", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
	//用户
	public_users_v1 := router.Group("/api/v1/users")
	{
		public_users_v1.POST("/register", sc.UserHandler.Register)
		public_users_v1.POST("/login", sc.UserHandler.Login)
	}

	//文章
	protected_post_v1 := router.Group("/api/v1/posts")
	protected_post_v1.Use(middleware.Auth([]byte(cfg.Jwt.Secret)))
	{
		protected_post_v1.GET("", sc.PostHandler.List)
		protected_post_v1.POST("", sc.PostHandler.PublishBlog)
		protected_post_v1.GET("/:id", sc.PostHandler.Get)
		protected_post_v1.PUT("/:id", sc.PostHandler.Update)
		protected_post_v1.DELETE("/:id", sc.PostHandler.Delete)
	}

	//评论
	protected_comment_v1 := router.Group("/api/v1/comments")
	protected_comment_v1.Use(middleware.Auth([]byte(cfg.Jwt.Secret)))
	{
		protected_comment_v1.POST("", sc.CommentHandler.Comment)
		protected_comment_v1.GET("", sc.CommentHandler.List)
		protected_comment_v1.DELETE("/:id", sc.CommentHandler.Delete)
	}
	return router
}

func initDb() *gorm.DB {
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(fmt.Errorf("创建数据目录失败: %w", err))
	}

	db, err := gorm.Open(sqlite.Open("data/blogs.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	// 根据 User 结构自动创建或更新表
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}); err != nil {
		panic(fmt.Errorf("创建数据表失败：: %w", err))
	}

	return db
}
