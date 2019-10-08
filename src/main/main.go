package main

import (
	"articlebk/src/utils/initinfo"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"io"
	"os"
	"time"
)

var (
	confPath = "../../conf/config.ini"
)

func main() {
	// 设置路由
	cfg, err := ini.Load(confPath)
	if err != nil {
		fmt.Println("initLog.go:22 无法加载配置文件内容.日志路径无法读取", err)
		os.Exit(1)
	}
	// 运行模式
	mode := cfg.Section("").Key("app_mode").String()
	router := gin.Default()
	gin.DisableConsoleColor()
	logfilePath := cfg.Section(mode).Key("logfile_path").String()
	// 创建记录日志的文件
	f, _ := os.Create(logfilePath)
	// 只输出日志到文件,不在终端打印
	gin.DefaultWriter = io.MultiWriter(f)
	// 启用gin的访问日志
	router.Use(gin.Logger())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	initinfo.InitLog()
	//mode := cfg.Section("").Key("app_mode").String()
	//httpPort := cfg.Section(mode).Key("server.address").String()
	//router := initinfo.InitLog(confPath)
	_ = router.Run(":8080")
}
