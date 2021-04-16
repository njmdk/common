package redisclient

import (
	"time"

	"github.com/njmdk/common/utils"
)

type Config struct {
	Addr       string         `json:"addr" toml:"addr"`
	Password   string         `json:"password" toml:"password"`
	DataBase   int64          `json:"data_base" toml:"data_base"`
	Timeout    utils.Duration `json:"timeout" toml:"timeout"`
	MaxRetries int64          `json:"max_retries" toml:"max_retries"`
}

func NewRedisConfig() *Config {
	return &Config{
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DataBase:   0,
		Timeout:    utils.Duration{Duration: time.Second},
		MaxRetries: 3,
	}
}
