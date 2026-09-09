package main

import (
	"homework01/homework04/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	//加载配置文件
	cfg := config.Load()

	//注册路由
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//启动服务
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	router.Run(addr) // listens on 0.0.0.0:8080 by default
}
