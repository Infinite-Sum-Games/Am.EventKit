package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/redis/go-redis/v9"
)

var Limiter *redis.Client

func InitRedis() (*redis.Client, error) {
	host := cmd.Env.RedisHost
	port := cmd.Env.RedisPort
	uname := cmd.Env.RedisUsername
	passwd := cmd.Env.RedisPassword
	resp := 3

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: uname,
		Password: passwd,
		DB:       0, // default DB
		Protocol: resp,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		pkg.Log.Fatal("[FAIL]: Health-check failed for Redis.", err)
		return nil, err
	}
	msg := fmt.Sprintf("[OK] Health-check successfuly for Redis. Response: %s", pong)
	pkg.Log.Info(msg)

	return rdb, nil
}

func CloseRedis(client *redis.Client) {
	if client != nil {
		client.Close()
	}
}
