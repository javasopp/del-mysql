package mysql

import (
	"fmt"
	"learngo/setting"
	"os/exec"
)

func DelFile() {
	// 执行外部 shell 脚本 "docker-compose-test/script.sh"
	cmd := exec.Command("/bin/sh", "-c", setting.AppConfig.MysqlInfo.File)

	// 捕获命令的输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// 打印输出结果
	fmt.Println("Command output:" + string(output))
}
