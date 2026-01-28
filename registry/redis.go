package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	reg := GetSingletonObject()
	if reg == nil {
		return nil, fmt.Errorf("registry not initialized")
	}

	config := reg.GetConfig()
	// config is a value, not a pointer, so we don't need nil check

	redisConfig := config.Redis
	host := redisConfig.Host
	port := redisConfig.Port
	uname := redisConfig.Username
	passwd := redisConfig.Password
	db := redisConfig.DB
	protocol := redisConfig.Protocol

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: uname,
		Password: passwd,
		DB:       db,
		Protocol: protocol,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	// Update registry with Redis client
	if clientSetter := reg; clientSetter != nil {
		// Since we can't directly modify the registry from here,
		// this will be handled in the initialization
	}

	return rdb, nil
}

func InitRedisWithOptions(host, port, username, password string, db, protocol int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: username,
		Password: password,
		DB:       db,
		Protocol: protocol,
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

// Legacy functions for backward compatibility
func getEnvRedisHost() string {
	reg := GetSingletonObject()
	if reg != nil {
		config := reg.GetConfig()
		return config.Redis.Host
	}
	return "localhost"
}

func getEnvRedisPort() string {
	reg := GetSingletonObject()
	if reg != nil {
		config := reg.GetConfig()
		return config.Redis.Port
	}
	return "6379"
}

func getEnvRedisUsername() string {
	reg := GetSingletonObject()
	if reg != nil {
		config := reg.GetConfig()
		return config.Redis.Username
	}
	return ""
}

func getEnvRedisPassword() string {
	reg := GetSingletonObject()
	if reg != nil {
		config := reg.GetConfig()
		return config.Redis.Password
	}
	return ""
}
