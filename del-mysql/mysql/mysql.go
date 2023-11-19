package mysql

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"learngo/setting"
	"learngo/utils"
	"strconv"
)

var connectTimes = 0
var continueTimes = 0
var err error

func Init() {
	// 从配置文件中获取 MySQL 连接信息
	dbUser := setting.AppConfig.MysqlInfo.User
	dbPassword := setting.AppConfig.MysqlInfo.Password
	dbHost := setting.AppConfig.MysqlInfo.Host
	dbPort := setting.AppConfig.MysqlInfo.Port
	dbName := setting.AppConfig.MysqlInfo.Database

	// 构建 MySQL 连接字符串
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	connectMysqlByTimes(dbSource)
	judgeTimes(dbSource)
}

func judgeTimes(dbSource string) {
	if continueTimes <= 0 {
		if connectTimes >= 4 {
			log.Info("mysql连接成功,i got you~")
			continueTimes = 1
		} else {
			log.Errorf("mysql连接错误，该容器出现异常~")
			log.Errorf("%v", err)
			log.Error("准备删除mysql数据卷了~")
			DelFile()
			log.Info("开始执行ping验证~")
			result := utils.Ping()
			if result == true {
				connectTimes = 0
				connectMysqlByTimes(dbSource)
				judgeTimes(dbSource)
			}
		}
	}
}

func connectMysqlByTimes(dbSource string) {
	for i := 0; i < 4; i++ {
		// 四次循环判断
		// 测试 MySQL 连通性
		err = testMySQLConnection(dbSource)
		log.Info("第" + strconv.Itoa(i+1) + "次连接~")
		if err == nil {
			connectTimes++
		}
	}
}

func testMySQLConnection(dataSource string) error {
	// 打开数据库连接
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// 尝试连接
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}
