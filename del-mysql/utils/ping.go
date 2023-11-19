package utils

import (
	"fmt"
	"learngo/setting"
	"os/exec"
	"strings"
	"time"
)

var pingTimes = 0

func Ping() bool {
	ipToPing := setting.AppConfig.MysqlInfo.Host
	for {
		// 执行 ping 命令
		pingCmd := exec.Command("ping", "-c", "5", ipToPing)
		pingOutput, err := pingCmd.CombinedOutput()

		if err != nil {
			fmt.Println("执行ping报错了:", err)
			fmt.Println(string(pingOutput))
			continue
		}

		// 检查 ping 结果
		if strings.Contains(string(pingOutput), "bytes from") {
			fmt.Println("ping成功了~")
			pingTimes++
			if pingTimes >= 4 {
				break
			}
		} else {
			fmt.Println("ping失败了，继续ping")
			time.Sleep(1 * time.Second)
		}
	}
	return true
}
