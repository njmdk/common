package redisclient

import (
	"time"

	"github.com/go-redis/redis"
)

// NewRedisClient 生成一个redis连接池client.
// 返回值：如果error==nil,success;否则的话,*redis.Client==nil
func NewRedisClient(
	addr string,
	db int64,
	password string,
	timeout time.Duration,
	maxRetries int64) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          int(db),
		DialTimeout: timeout,
		MaxRetries:  int(maxRetries),
		PoolSize:    50,
	})

	err := rc.Ping().Err()
	if err != nil {
		return nil, err
	}

	return rc, nil
}

// NewRedisClient 生成一个redis连接池client.
// 返回值：如果error==nil,success;否则的话,*redis.Client==nil
func NewRedisClientWithConfig(cfg *Config) (*redis.Client, error) {
	return NewRedisClient(cfg.Addr, cfg.DataBase, cfg.Password, cfg.Timeout.Duration, cfg.MaxRetries)
}

var globalRedis *redis.Client
var _ = globalRedis

func NewAndSetDefaultRedis(cfg *Config) error {
	var err error
	globalRedis, err = NewRedisClientWithConfig(cfg)

	return err
}
