package db

import (
	"ray_box/infrastructure/config"

	"github.com/redis/go-redis/v9"
)

const (
	REDIS_DB_MASTER = "master"
)

var (
	redisConn = make(map[string]redis.UniversalClient)
	// lock      = sync.RWMutex{}
)

func GetRedisConn(name string) redis.UniversalClient {
	// lock.RLock()
	conn, ok := redisConn[name]
	// lock.RUnlock()

	if !ok {
		// lock.Lock()
		options := config.RedisConfig(name)
		redisConn[name] = redis.NewUniversalClient(options)
		conn = redisConn[name]
		// lock.Unlock()
	}

	return conn
}
