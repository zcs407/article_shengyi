package main

import (
	"articlebk/src/common"
	"articlebk/src/common/database"
	router2 "articlebk/src/common/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

var (
	confPath = "conf/config.toml"
)

func main() {

	common.InitConfig(confPath, common.Settings)
	common.InitLog(common.Settings.Logging.Path, common.Settings.Logging.Level, common.Settings.Logging.Format)
	database.InitDB(common.Settings.Database)
	common.Log.Info("初始化完毕,启动api服务,配置文件路径为:%s", confPath)
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
	router = router2.InitRouter()
	common.Log.Info("访问地址:http://%s", common.Settings.ApiServer.Address)
	_ = router.Run(common.Settings.ApiServer.Address)
}
