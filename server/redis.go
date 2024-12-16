package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var redisManager redisWrapper

type redisBuffer interface {
	saveToBuff(chatName, sender string)
}

type redisWrapper struct {
	redis *redis.Client
}

func handleOfflineMessages(cl *Client) {
	data := map[string][]interface{}{}
	for _, sub := range cl.subs {
		key := fmt.Sprintf("offline:messages:%s:%s", sub, cl.username)
		fmt.Println("key", key)
		if ok := redisManager.hasMessages(key); ok {
			fmt.Println("exist", ok)
			data[sub] = redisManager.getOffMessages(key)
		}
	}
	if len(data) > 0 {
		writeToSocket(data, "OFFLINE_MESSAGES", cl, websocket.TextMessage)
	}
}

func (r *redisWrapper) hasMessages(key string) bool {
	var ctx = context.Background()
	len, err := r.redis.LLen(ctx, key).Result()
	fmt.Println("len", len)
	if err != nil {
		fmt.Println("Error checking for messages", err)
	}
	if len > 0 {
		return true
	}
	fmt.Println("No messages to retrieve")
	return false
}

func (r *redisWrapper) getOffMessages(key string) []interface{} {
	var ctx = context.Background()
	messages := []interface{}{}
	req := r.redis.LRange(ctx, key, 0, -1)
	strSlice, _ := req.Result()

	for _, val := range strSlice {
		buffMessage := []byte(val)
		var redisContainer JSONenvelope

		var err = json.Unmarshal(buffMessage, &redisContainer)
		data := kindHandlers[redisContainer.Type]
		json.Unmarshal(buffMessage, &data)
		messages = append(messages, data)
		if err != nil {
			log.Fatal("error scanning slice", err)
		}
	}

	fmt.Println("res messgs", messages)
	for i := 0; i < len(messages); i++ {
		r.redis.LPop(ctx, key)
	}
	fmt.Println("redis messages:", messages)
	return messages
}

func (r *redisWrapper) offlineMessageProcessing(msg interface{}, msgType, chatName, sender string) {
	var b redisBuffer
	switch msgType {
	case "NEW_MESSAGE":
		data := msg.(UserMessage)
		b = data

	case "NEW_PRIVATE_CHAT":
		// data := msg.(NewPrivateChat)
		// b = data
	}
	b.saveToBuff(chatName, sender)
}

func (um UserMessage) saveToBuff(chatName, sender string) {
	var ctx = context.Background()
	var receivers []string
	if chatType(chatName) == "private" {
		subManager := dbManager.initializeDBhandler("subscription")
		receivers = subManager.GetPrivateChatSubs(chatName, sender)
	} else {
		// ...
	}
	var redisContainer struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}
	redisContainer.Type = "NEW_MESSAGE"
	redisContainer.Data = um
	json, err := json.Marshal(redisContainer)
	if err != nil {
		fmt.Println("Error while serializing cache data to redis")
	}

	for _, receiver := range receivers {
		key := fmt.Sprintf("offline:messages:%s:%s", chatName, receiver)
		fmt.Println("key save", key)
		redisManager.redis.RPush(ctx, key, json)
		if err != nil {
			fmt.Println("error retrieving redis message", err)
		}
	}
}

func chatType(chatName string) string {
	if len(strings.Split(chatName, "_")) > 1 {
		return "private"
	} else {
		return "group"
	}
}
