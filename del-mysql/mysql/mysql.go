package mysql

import (
	"bufio"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"learngo/setting"
	"learngo/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	isDel := setting.AppConfig.MysqlInfo.IsDel
	if isDel == 1 {
		connectMysqlByTimes(dbSource)
		judgeTimes(dbSource)
	}
}

func BackupDatabase() error {
	fileName := fmt.Sprintf("%s-new.sql", setting.AppConfig.MysqlInfo.Database)
	tempFileName := fmt.Sprintf(".%s.new.sql.tmp", setting.AppConfig.MysqlInfo.Database)
	tempFilePath := filepath.Join("/var/lib/mysql/", tempFileName)
	finalFilePath := filepath.Join("/var/lib/mysql/", fileName)

	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer tempFile.Close()

	cmdArgs := []string{
		"-h", setting.AppConfig.MysqlInfo.Host,
		"-P", setting.AppConfig.MysqlInfo.Port,
		"-u", setting.AppConfig.MysqlInfo.User,
		"-p" + setting.AppConfig.MysqlInfo.Password,
		"--single-transaction",
		setting.AppConfig.MysqlInfo.Database,
	}

	cmdStr := fmt.Sprintf("mysqldump %s", strings.Join(cmdArgs, " "))
	log.Info("执行的 mysqldump 命令: ", cmdStr)

	cmd := exec.Command("mysqldump", cmdArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取标准输出管道失败: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取错误输出管道失败: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动mysqldump失败: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Error("读取mysqldump标准输出时出错:", err)
				}
				break
			}
			log.Info("mysqldump标准输出:", line)
			if _, err := tempFile.WriteString(line); err != nil {
				log.Error("写入临时文件失败:", err)
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Error("读取mysqldump错误输出时出错:", err)
				}
				break
			}
			log.Error("mysqldump错误输出:", line)
		}
	}()

	// 等待所有的 stdout 和 stderr 读取完成
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("执行mysqldump失败: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("同步临时文件失败: %w", err)
	}

	if err := os.Rename(tempFilePath, finalFilePath); err != nil {
		return fmt.Errorf("重命名临时文件失败: %w", err)
	}

	return nil
}

func judgeTimes(dbSource string) {
	if continueTimes <= 0 {
		if connectTimes >= 4 {
			log.Info("mysql连接成功,i got you~")
			continueTimes = 1
			if setting.AppConfig.MysqlInfo.IsBox == 1 {
				log.Info("开启之前的备份")
				_ = BackupDatabase()
			}
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
