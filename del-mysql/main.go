package main

import (
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
}
