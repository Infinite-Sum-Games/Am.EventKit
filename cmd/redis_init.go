package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	redis "github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func InitCache() error {
	host := Env.RedisHost
	port := Env.RedisPort
	uname := Env.RedisUsername
	passwd := Env.RedisPassword
	resp := 3

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: uname,
		Password: passwd,
		DB:       0,
		Protocol: resp,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		pkg.Log.LogFatal("[FAIL]: Heal-check failed for Redis", err)
		return err
	}
	pkg.Log.LogInfo(
		fmt.Sprintf("[SUCCESS]: Health-check successful for Redis. Response: %s", pong),
	)

	Cache = rdb
	return nil
}
