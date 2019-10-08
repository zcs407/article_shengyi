package main

import (
	"articlebk/src/common"
	router2 "articlebk/src/common/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
)

var (
	confPath = "/Users/ander/go/src/articlebk/conf/config.toml"
)

func main() {

	router := gin.Default()
	gin.DisableConsoleColor()

	// 创建记录日志的文件
	f, _ := os.Create(common.Settings.Logging.Path)
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
	log.Printf("配置文件路径为:%s\n", confPath)
	common.InitConfig(confPath, common.Settings)
	common.InitLog(common.Settings.Logging.Path, common.Settings.Logging.Level, common.Settings.Logging.Format)

	common.InitDB(common.Settings.Database)
	log.Println("目前没问题")

	router = router2.InitRouter()
	_ = router.Run(common.Settings.ApiServer.Address)
}
