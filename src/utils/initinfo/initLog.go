package initinfo

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"io"
	"log"
	"os"
)

var (
	INFO       *log.Logger
	WARING     *log.Logger
	ERROR      *log.Logger
	DBERR      *log.Logger
	configFile = "../../conf/config.ini"
)

func InitLog() {
	// 加载配置
	cfg, err := ini.Load(configFile)
	if err != nil {
		log.Println("initLog.go:22 无法加载配置文件内容.日志路径无法读取", err)
		os.Exit(1)
	}
	// 运行模式
	mode := cfg.Section("").Key("app_mode").String()
	// 日志路径获取
	logPath := cfg.Section(mode).Key("paths.log").String()
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("initLog.go:31  无法创建爬虫日志文件,请检查路径", err)
		return
	}
	INFO = log.New(gin.DefaultWriter, "[INFO]", log.LstdFlags|log.Lshortfile)
	WARING = log.New(io.MultiWriter(file), "[WARING]", log.LstdFlags|log.Lshortfile)
	ERROR = log.New(io.MultiWriter(file), "[ERROR]", log.LstdFlags|log.Lshortfile)
	DBERR = log.New(io.MultiWriter(file), "[DB_ERR]", log.LstdFlags|log.Lshortfile)
	// 设置运行模式
	mode = cfg.Section("").Key("app_mode").String()
	if mode == "develop" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
