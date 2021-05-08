package utils

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
)

func NewCaptcha() *base64Captcha.Captcha {
	var store = base64Captcha.DefaultMemStore
	//var driver = base64Captcha.DefaultDriverDigit
	var driver = base64Captcha.NewDriverDigit(50,100,4,0.7,80)

	return base64Captcha.NewCaptcha(driver,store)
}

type RedisStore struct {
	redisClient redis.Cmdable
}

func (this_ *RedisStore) Set(id string, value string) {
	this_.redisClient.Set("captcha:"+id, value, time.Minute*5)
}

func (this_ *RedisStore) Get(id string, clear bool) string {
	str := this_.redisClient.Get("captcha:" + id).Val()
	if clear && str != "" {
		this_.redisClient.Del("captcha:" + id)
	}

	return str
}
