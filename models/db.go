package models

import "github.com/go-redis/redis"

var client *redis.Client

// Init is uset to initialize a redis client
func Init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
