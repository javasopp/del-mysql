package setting

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	MysqlInfo Mysql `mapstructure:"mysql"`
}

type Mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"mysql"`
	File     string `yaml:"file"`
	IsDel    int    `yaml:"is_del"`
}

var AppConfig Config

func init() {
	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".") // 指定配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}
}
