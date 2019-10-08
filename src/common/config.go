package common

import (
	"errors"
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
)

var Settings *Config = &Config{}

type Logging struct {
	Level  string
	Format string
	Path   string
}
type Database struct {
	Address      string
	UserName     string `toml:"user_name"`
	Password     string
	DbName       string `toml:"db_name"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `toml:"max_open_conns"`
}
type ApiServer struct {
	Address string
}
type FileServer struct {
	Address string
}
type token struct {
}
type Config struct {
	Logging    `toml:"logging"`
	Database   `toml:"database"`
	FileServer `toml:"file_server"`
	ApiServer  `toml:"api_server"`
}

func InitConfig(path string, config interface{}) {
	if len(path) == 0 {
		panic(errors.New("没有提供配置文件"))
	}
	configPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	if _, err = toml.DecodeFile(configPath, config); err != nil {
		panic(err)
	} else {
		log.Println("配置文件路径:", path)
	}
}
