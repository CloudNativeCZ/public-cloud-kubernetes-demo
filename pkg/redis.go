package pkg

import (
	"github.com/go-redis/redis"
)

func NewClient(options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)

	_, err := client.Ping().Result()
	if err != nil {
		return client, err
	}

	return client, nil
}