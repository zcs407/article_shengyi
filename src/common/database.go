package common

import (
	"articlebk/src/common/dbtable"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var DBSQL *gorm.DB

func InitDB(mysql Database) {
	//Customize the datetime
	gorm.NowFunc = func() time.Time {
		return time.Now().Round(time.Second)
	}
	dsn := mysql.UserName + ":" + mysql.Password + "@tcp(" + mysql.Address + ")/" + mysql.DbName + "?charset=utf8&parseTime=true&loc=Local"
	Log.Info("连接数据库 %s", dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		Log.Error("无法连接数据库", err)
		panic(err)
	}
	db.DB()
	db.LogMode(true)

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	//设置最大空闲连接池
	db.DB().SetMaxIdleConns(mysql.MaxIdleConns)
	//设置最大打开连接池
	db.DB().SetMaxOpenConns(mysql.MaxOpenConns)
	//连接最大超时时间
	db.DB().SetConnMaxLifetime(time.Hour)
	db.SingularTable(true)
	db.AutoMigrate(&dbtable.User{}, &dbtable.Article{}, &dbtable.Image{}, &dbtable.Special{}, &dbtable.Tag{})
	go func() {
		timer := time.NewTicker(5 * time.Second)
		for {
			if _, ok := <-timer.C; !ok {
				continue
			}
			Log.Warn("当前打开数据库连接数为: %d", db.DB().Stats().OpenConnections)
		}
	}()
	DBSQL = db
}
