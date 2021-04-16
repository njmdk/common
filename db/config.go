package db

import (
	"time"

	"github.com/njmdk/common/utils"
)

type Config struct {
	Addr         string         `toml:"addr" json:"addr"`
	User         string         `toml:"user" json:"user"`
	Password     string         `toml:"password" json:"password"`
	Database     string         `toml:"database" json:"database"`
	MaxOpenConnS int            `toml:"max_open_conn_num" json:"max_open_conn_s"`
	MaxIdleConnS int            `toml:"max_idle_conn_num" json:"max_idle_conn_s"`
	MaxLifeTime  utils.Duration `toml:"max_life_time" json:"max_life_time"`
}

func NewDefaultMysqlConfig() *Config {
	return &Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Password:     "root",
		Database:     "auth",
		MaxOpenConnS: 300,
		MaxIdleConnS: 1,
		MaxLifeTime:  utils.Duration{Duration: time.Second * 300},
	}
}
