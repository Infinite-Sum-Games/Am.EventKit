package cmd

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var Valkey *redis.Client

func InitValkey() (*redis.Client, error) {
	host := Env.RedisHost
	port := Env.RedisPort
	resp := 3

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0, // default DB
		Protocol: resp,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result() // health-check
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func CloseValkey(client *redis.Client) {
	if client != nil {
		client.Close()
	}
}
