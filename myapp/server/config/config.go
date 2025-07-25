package config

import (
	"fmt"
	"os"
	"strconv"
	"github.com/spf13/viper"
)

type Config struct {
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
	}
}

var Appconf Config

func Inits() {
	// 优先使用环境变量
	if os.Getenv("DB_HOST") != "" {
		Appconf.Mysql.Host = os.Getenv("DB_HOST")
		Appconf.Mysql.User = getEnvOrDefault("DB_USER", "root")
		Appconf.Mysql.Password = getEnvOrDefault("DB_PASSWORD", "root123456")
		Appconf.Mysql.Database = getEnvOrDefault("DB_NAME", "myapp")
		
		portStr := getEnvOrDefault("DB_PORT", "3306")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			port = 3306
		}
		Appconf.Mysql.Port = port
		
		fmt.Println("使用环境变量配置数据库连接")
		return
	}
	
	// 回退到配置文件
	viper.SetConfigFile("./config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		fmt.Println("配置文件不存在，使用默认配置")
		Appconf.Mysql.Host = "localhost"
		Appconf.Mysql.User = "root"
		Appconf.Mysql.Password = "root123456"
		Appconf.Mysql.Database = "myapp"
		Appconf.Mysql.Port = 3306
		return
	}
	
	err = viper.Unmarshal(&Appconf)
	if err != nil {
		panic(err)
	}
	fmt.Println("配置文件加载成功")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
