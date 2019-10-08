package initinfo

import (
	"articlebk/src/utils/dbtable"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

var (
	DB         *gorm.DB
	configPath = "../../conf/config.ini"
)

func init() {
	log.Println("初始化DB成功")
	// 加载配置
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		log.Println("数据库加载配置文件错误,无法读取配置信息", err)
		os.Exit(10000)
	}
	// 运行模式
	mode := cfg.Section("").Key("app_mode").String()
	// 主机
	host := cfg.Section(mode).Key("mysql.host").String()
	// 端口
	port := cfg.Section(mode).Key("mysql.port").String()
	// 用户名
	username := cfg.Section(mode).Key("mysql.username").String()
	// 密码
	password := cfg.Section(mode).Key("mysql.password").String()
	// 数据库名称
	dbname := cfg.Section(mode).Key("mysql.dbname").String()
	// 最大空闲连接数
	maxIdleConns, err := cfg.Section(mode).Key("mysql.max_idle_conns").Int()
	if err != nil {
		fmt.Println(maxIdleConns)
		fmt.Printf("%v", err)
		os.Exit(100000)
	}
	// 最大打开的连接数
	maxOpenConns, err := cfg.Section(mode).Key("mysql.max_open_conns").Int()
	if err != nil {
		log.Printf("%v", err)
		os.Exit(100000)
	}

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=true&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println("无法连接数据库", err)
		os.Exit(10000)
	}

	//设置最大空闲连接池
	db.DB().SetMaxIdleConns(maxIdleConns)
	//设置最大打开连接池
	db.DB().SetMaxOpenConns(maxOpenConns)
	//连接最大超时时间
	db.DB().SetConnMaxLifetime(time.Hour)
	db.SingularTable(true)
	db.AutoMigrate(&dbtable.User{}, &dbtable.Article{}, &dbtable.Image{}, &dbtable.Special{}, &dbtable.Tag{})
	DB = db
}
