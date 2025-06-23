package sessionstore

import (
	"push/common/lib/env"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(env env.Env) *redis.Client {
	redisAddr := "localhost:" + env.Redis.Port
	return redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     "", // TODO  : redis 비밀번호 설정 (암호화 필요)
		DB:           0,
		DialTimeout:  2 * time.Second, // tcp 연결 타임아웃시간
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     10,
	})
}
