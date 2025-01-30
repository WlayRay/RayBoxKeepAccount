package db_test

import (
	"context"
	"fmt"
	"ray_box/infrastructure/db"
	"testing"
	"time"
)

func TestRedisClient(t *testing.T) {
	redisConn := db.GetRedisConn(db.REDIS_DB_MASTER)
	defer redisConn.Close()

	timeoutContext, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := redisConn.Ping(timeoutContext).Result()
	if err != nil {
		t.Errorf("redis ping error: %v", err)
	} else {
		fmt.Println("redis ping response:", pong)
	}
}
