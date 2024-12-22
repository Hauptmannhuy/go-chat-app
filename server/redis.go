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
	saveToBuff()
}

type redisContainer struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
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
		var redisTypeContainer JSONenvelope
		var dataContainer redisContainer
		var err = json.Unmarshal(buffMessage, &redisTypeContainer)
		data := kindHandlers[redisTypeContainer.Type]
		json.Unmarshal(buffMessage, &data)
		dataContainer.Type = redisTypeContainer.Type.toValue()
		dataContainer.Data = data
		messages = append(messages, dataContainer)
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

func (r *redisWrapper) offlineMessageProcessing(msg interface{}, msgType string) {
	var b redisBuffer

	switch msgType {
	case "NEW_MESSAGE":
		data := msg.(UserMessage)
		b = data

	case "NEW_PRIVATE_CHAT":
		data := msg.(NewPrivateChat)
		b = data

	}
	b.saveToBuff()
}

func (um UserMessage) saveToBuff() {
	var ctx = context.Background()
	var receivers []string
	if chatType(um.ChatName) == "private" {
		subManager := dbManager.initializeDBhandler("subscription")
		receivers = subManager.GetPrivateChatSubs(um.ChatName, um.Username)
	} else {
		// ...
	}
	container := redisContainer{
		Type: "NEW_MESSAGE",
		Data: um,
	}

	json, err := json.Marshal(container)
	if err != nil {
		log.Println("Error while serializing cache data to redis")
	}

	for _, receiver := range receivers {
		key := fmt.Sprintf("offline:messages:%s:%s", um.ChatName, receiver)
		fmt.Println("key save", key)
		redisManager.redis.RPush(ctx, key, json)
		if err != nil {
			fmt.Println("error retrieving redis message", err)
		}
	}
}

func (npc NewPrivateChat) saveToBuff() {
	ctx := context.Background()
	receiverName := strings.Split(npc.ChatName, "_")[1]
	key := fmt.Sprintf("offline:messages:%s:%s", npc.ChatName, receiverName)
	container := redisContainer{
		Type: "NEW_PRIVATE_CHAT",
		Data: npc,
	}
	json, err := json.Marshal(container)
	if err != nil {
		log.Println(err)
	}
	redisManager.redis.LPush(ctx, key, json)
}

func chatType(chatName string) string {
	if len(strings.Split(chatName, "_")) > 1 {
		return "private"
	} else {
		return "group"
	}
}
