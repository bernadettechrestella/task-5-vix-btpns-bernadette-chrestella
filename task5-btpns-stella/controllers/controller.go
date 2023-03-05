package controllers

import (
	storage "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"strings"
	"task5-btpns-stella/middlewares"
	"time"
)

var (
	client = gocloak.NewClient("http://localhost:8080")
	ctx    = context.Background()
)

func InitFirebaseStorage() (*storage.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	opt := option.WithCredentialsFile("./certs/serviceAccountKey.json")
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return client, nil
}

func authorized(c *gin.Context) (string, *middlewares.RedisValue, error) {
	beartoken := c.Request.Header.Get("Authorization")
	if beartoken == "" {
		return "", nil, errors.New("Token Not Found")
	}

	tokenExtract := strings.Split(beartoken, " ")
	if len(tokenExtract) != 2 {
		return "", nil, errors.New("Cannot Split Token")
	}

	redisValue, err := middlewares.ReadRedis(tokenExtract[1])
	if err != nil {
		return "", nil, err
	}

	var payloadToken *middlewares.RedisValue

	err = json.Unmarshal([]byte(redisValue), &payloadToken)
	if err != nil {
		return "", nil, err
	}

	return tokenExtract[1], payloadToken, nil
}
