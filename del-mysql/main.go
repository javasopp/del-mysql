package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"learngo/mysql"
	"time"
)

func main() {
	log.Info("等待10s，等待mysql容器构建~")
	time.Sleep(10 * time.Second)
	log.Info("10s过后了，开始执行逻辑判断~")
	mysql.Init()

	// 初始化gin
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 定义保存数据库的表
	r.GET("/save", func(c *gin.Context) {
		path, err := mysql.BackupDatabase()
		if err != nil {
			panic(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"backup_file": path,
			"code":        0,
		})
	})
	r.Run(":8090") // 监听并在 0.0.0.0:8080 上启动服务
}
