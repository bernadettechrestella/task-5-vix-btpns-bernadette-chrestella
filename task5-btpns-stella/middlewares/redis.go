package middlewares

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var (
	RedisHost = os.Getenv("URL_CACHE_DB_HOST")
	RedisPort = os.Getenv("URL_CACHE_DB_PORT")
)

func contextRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisHost + ":" + RedisPort,
		Username: "user-account",
		Password: "cWQTp6He@5scCWf",
		DB:       0,
	})
	return client
}

func WriteRedis(key string, body string, expiry time.Duration) (string, error) {
	client := contextRedis()
	status, err := client.Set(context.Background(), key, body, expiry).Result()
	defer client.Close()
	return status, err
}

var ReadRedis = func(key string) (string, error) {
	client := contextRedis()
	resultMessage, err := client.Get(context.Background(), key).Result()
	defer client.Close()
	return resultMessage, err
}

var PurgeRedis = func(key string) (int64, bool) {
	var errResult = false
	client := contextRedis()
	status, _ := client.Del(context.Background(), key).Result()
	if status == 1 {
		errResult = false
	} else if status == 0 {
		errResult = true
	}
	defer client.Close()
	return status, errResult
}

type RedisValue struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}
