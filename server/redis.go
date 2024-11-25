package main

import "github.com/redis/go-redis/v9"

var redisManager redisDBwrap

type redisDBwrap struct {
	redis *redis.Client
}
