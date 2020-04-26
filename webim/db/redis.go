package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var (
	defaultRedisReadTimeout  = 2 * time.Second
	defaultRedisWriteTimeout = 2 * time.Second
)

// RedisClient redis client
var RedisClient *redis.Client

// NewRedisClient create redis client
func NewRedisClient(redisConf *conf.RedisConfig) error {
	if redisConf.Addr == "" {
		return fmt.Errorf("redis addr empty")
	}

	readTimeout := defaultRedisReadTimeout
	if redisConf.ReadTimeout.Duration > 0 {
		readTimeout = redisConf.ReadTimeout.Duration
	}

	writeTimeout := defaultRedisWriteTimeout
	if redisConf.WriteTimeout.Duration > 0 {
		writeTimeout = redisConf.WriteTimeout.Duration
	}

	RedisClient = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         redisConf.Addr,
		DB:           0,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	_, err := RedisClient.Ping().Result()
	return err
}
