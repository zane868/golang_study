package main

import (
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/zane868/golang_study/homework04/config"
	"github.com/zane868/golang_study/homework04/handler"
	"github.com/zane868/golang_study/homework04/logging"
	"github.com/zane868/golang_study/homework04/middleware"
	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

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
	logger, logFile, err := logging.New("logs/app.log", os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)
	// Gin 的调试信息也经由 slog 同时写入控制台和文件。
	gin.DebugPrintFunc = func(format string, values ...interface{}) {
		slog.Info(fmt.Sprintf(format, values...))
	}
	exitCode := 0
	if err := run(); err != nil {
		slog.Error("Application startup failed", "error", err)
		exitCode = 1
	}
	if err := logFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "关闭日志文件失败: %v\n", err)
		exitCode = 1
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func run() error {

	//加载配置文件
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.Server.Mode != gin.DebugMode && cfg.Server.Mode != gin.ReleaseMode && cfg.Server.Mode != gin.TestMode {
		return fmt.Errorf("invalid server mode")
	}
	gin.SetMode(cfg.Server.Mode)
	slog.Info("Configuration loaded", "mode", cfg.Server.Mode)

	//初始化数据库
	db, err := initDb()
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("Database close failed", "error", err)
		}
	}()
	slog.Info("Database initialized", "driver", "sqlite")

	//实例化服务
	serviceContext := NewContext(db, cfg)

	//注册路由
	router := regRouter(serviceContext, cfg)

	//启动服务
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	slog.Info("HTTP server starting", "address", addr)
	return router.Run(addr)
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
	router := gin.New()
	router.Use(middleware.RequestLog(), middleware.Recover())
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

func initDb() (*gorm.DB, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	// 数据库错误由统一错误处理记录，避免默认 SQL 日志输出密码哈希等参数。
	db, err := gorm.Open(sqlite.Open("data/blogs.db"), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 根据 User 结构自动创建或更新表
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}); err != nil {
		if sqlDB, closeErr := db.DB(); closeErr == nil {
			sqlDB.Close()
		}
		return nil, fmt.Errorf("创建数据表失败: %w", err)
	}

	return db, nil
}
