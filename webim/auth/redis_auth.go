package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var (
	defaultAuthInfoExpire = 24 * time.Hour
)

type RedisAgentAuth struct {
	redisCli redis.Cmdable
}

func NewRedisAgentAuth(pipeliner redis.Cmdable) *RedisAgentAuth {
	return &RedisAgentAuth{
		redisCli: pipeliner,
	}
}

func TokenExists(redisCli redis.Cmdable, token string) (bool, error) {
	v, err := redisCli.Exists(token).Result()
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

func RemoveToken(redisCli redis.Cmdable, token string) error {
	return redisCli.Del(token).Err()
}

func SignOutAll(redisCli redis.Cmdable, agentToken, entID, agentID string) (err error) {
	agentTokensKey := fmt.Sprintf(common.AgentTokenList, agentID)
	tokens, err := redisCli.SMembers(agentTokensKey).Result()
	if err != nil {
		return err
	}

	pipe := redisCli.Pipeline()
	defer func() {
		if closeErr := pipe.Close(); closeErr != nil {
			log.Logger.Warnf("close pipe error: %v", closeErr)
		}
	}()

	var del bool
	var userTokens []interface{}
	for _, token := range tokens {
		userTokens = append(userTokens, token)

		if token == agentToken {
			del = true
		}

		pipe.Del(token)
	}

	if !del {
		pipe.Del(agentToken)
	}

	pipe.Del(fmt.Sprintf(common.AgentLoginCount, agentID))
	pipe.Del(agentTokensKey)

	pipe.ZRem(fmt.Sprintf(common.EntOnlineAgentList, entID), userTokens...)

	_, err = pipe.Exec()
	return err
}

func (rauth *RedisAgentAuth) RegisterLogin(agentID, token string) (err error) {
	pipe := rauth.redisCli.Pipeline()
	defer pipe.Close()

	key := fmt.Sprintf(common.AgentTokenList, agentID)
	pipe.SAdd(key, token)

	key = fmt.Sprintf(common.AgentLoginCount, agentID)
	pipe.Incr(key)

	markAgentToken(pipe, token)
	_, err = pipe.Exec()
	return err
}

func (a *RedisAgentAuth) IncrLoginCount(pipe redis.Pipeliner, agentID string) error {
	key := fmt.Sprintf(common.AgentLoginCount, agentID)
	return pipe.Process(pipe.Incr(key))
}

func (a *RedisAgentAuth) DecrLoginCount(agentID string) error {
	key := fmt.Sprintf(common.AgentLoginCount, agentID)
	return a.redisCli.Decr(key).Err()
}

func (a *RedisAgentAuth) AppendAgentToken(agentID string, token string, pipe redis.Pipeliner) error {
	key := fmt.Sprintf(common.AgentTokenList, agentID)
	return pipe.Process(pipe.SAdd(key, token))
}

func (a *RedisAgentAuth) RemTokenFromList(agentID, token string) error {
	key := fmt.Sprintf(common.AgentTokenList, agentID)
	return a.redisCli.SRem(key, token).Err()
}

func (a *RedisAgentAuth) RemAllTokens(agentID string) error {
	key := fmt.Sprintf(common.AgentTokenList, agentID)
	return a.redisCli.Del(key).Err()
}

func (a *RedisAgentAuth) GetLoginCount(agentID string) (int, error) {
	key := fmt.Sprintf(common.AgentLoginCount, agentID)
	v, err := a.redisCli.Get(key).Result()
	if err != nil {
		return -1, err
	}

	loginCount, err := strconv.Atoi(v)
	if err != nil {
		return -1, err
	}

	return loginCount, nil
}

func markAgentToken(pipe redis.Pipeliner, token string) {
	var expire = defaultAuthInfoExpire
	var confD = conf.IMConf.AgentConf.AgentTokenExpire.Duration
	if confD > 0 && confD <= 10*24*time.Hour {
		expire = confD
	}
	pipe.Set(token, 1, expire)
}
