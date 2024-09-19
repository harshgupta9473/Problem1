package components

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func InitRedis()(*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB:0,
	})
	ping,err:=client.Ping(context.Background()).Result()
	if err!=nil{
		return nil,err
	}

	fmt.Println(ping)

	return client,nil
}

