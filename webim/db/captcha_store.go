package db

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
)

type RedisStore struct {
	Expire time.Duration
}

func (rs *RedisStore) Set(id string, digits []byte) {
	RedisClient.Set(id, string(digits), rs.Expire)
}

func (rs *RedisStore) Get(id string, clear bool) (digits []byte) {
	bs, err := RedisClient.Get(id).Bytes()
	if err != nil {
		log.Logger.Warnf("Get %s error: %v", id, err)
	}

	return bs
}
