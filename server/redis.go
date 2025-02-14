package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisDB redisManager

type redisBuffer interface {
	saveToBuff(receiverID int)
}

type redisContainer struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type redisManager struct {
	redis *redis.Client
}

func getRedis() *redisManager {
	if redisDB.redis == nil {
		fmt.Println("2")

		redisDB = redisManager{
			redis: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
				Protocol: 3,
			}),
		}
	}
	return &redisDB
}

func (r *redisManager) hasMessages(key string) bool {
	var ctx = context.Background()
	len, err := r.redis.LLen(ctx, key).Result()
	if err != nil {
		fmt.Println("Error checking for messages", err)
	}
	if len > 0 {
		return true
	}
	return false
}

func (r *redisManager) getOffMessages(key string) []interface{} {
	var ctx = context.Background()
	messages := []interface{}{}
	req := r.redis.LRange(ctx, key, 0, -1)
	strSlice, _ := req.Result()

	for _, val := range strSlice {
		buffMessage := []byte(val)
		var importRedisContainer struct {
			Type Kind            `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		var resultData redisContainer

		err := json.Unmarshal(buffMessage, &importRedisContainer)
		if err != nil {
			fmt.Println("error type unmarshaling redis message", err)
		}

		data, err := kindHandler(importRedisContainer.Type)
		resultData.Data = data

		if err != nil {
			fmt.Println("error kind handler", err)
		}

		err = json.Unmarshal(importRedisContainer.Data, &resultData.Data)

		if err != nil {
			fmt.Println("error data unmarshaling redis message", err)
		}

		resultData.Type = importRedisContainer.Type.toValue()
		resultData.Data = data
		messages = append(messages, resultData)

	}

	for i := 0; i < len(messages); i++ {
		r.redis.LPop(ctx, key)
	}
	return messages
}

func (um UserMessage) saveToBuff(receiverID int) {
	var ctx = context.Background()

	container := redisContainer{
		Type: "NEW_MESSAGE",
		Data: um,
	}

	json, err := json.Marshal(container)
	if err != nil {
		log.Println("Error while serializing cache data to redis")
	}

	key := fmt.Sprintf("offline:messages:%s:%d", um.ChatName, receiverID)
	getRedis().redis.RPush(ctx, key, json)

	if err != nil {
		fmt.Println("error retrieving redis message", err)
	}
}

func (npc NewPrivateChat) saveToBuff(receiverID int) {
	ctx := context.Background()
	key := fmt.Sprintf("offline:messages:%s:%d", npc.ChatName, receiverID)

	container := redisContainer{
		Type: "NEW_PRIVATE_CHAT",
		Data: npc,
	}

	json, err := json.Marshal(container)
	if err != nil {
		log.Println(err)
	}

	getRedis().redis.LPush(ctx, key, json)
}
