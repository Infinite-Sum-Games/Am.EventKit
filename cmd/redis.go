package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis() (*redis.Client, error) {
	host := Env.RedisHost
	port := Env.RedisPort
	uname := Env.RedisUsername
	passwd := Env.RedisPassword
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

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func CloseRedis(client *redis.Client) error {
	if client != nil {
		if err := client.Close(); err != nil {
			return err
		}
	}
	return nil

}
