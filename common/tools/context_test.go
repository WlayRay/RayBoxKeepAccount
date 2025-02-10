package tools_test

import (
	"ray_box/common/tools"
	"ray_box/infrastructure/db"
	"sync"
	"testing"
	"time"
)

func TestTimeoutCtx(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		{
			defer wg.Done()
			timeoutContext, cancel := tools.GetTimeoutCtx("redis")
			defer cancel()

			redisConn := db.GetRedisConn(db.REDIS_DB_MASTER)
			res, err := redisConn.Ping(timeoutContext).Result()
			if err != nil {
				t.Error(err)
			} else {
				t.Log(res)
			}

			time.Sleep(time.Second * 2)
			res, err = redisConn.Ping(timeoutContext).Result()
			if err != nil {
				t.Error(err)
			} else {
				t.Log(res)
			}
		}
	}()

	redisConn := db.GetRedisConn(db.REDIS_DB_MASTER)
	timeoutContext, cancel := tools.GetTimeoutCtx("redis")
	defer cancel()
	res, err := redisConn.Ping(timeoutContext).Result()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
	wg.Wait()
}
