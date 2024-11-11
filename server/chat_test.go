package main

import (
	"sync"
	"testing"
)

func TestCreateChat(t *testing.T) {
	var chatList ChatList
	chatList.Chats = make(map[string]*Chat)
	name := "chat room"
	t.Run("Chat creation ", func(t *testing.T) {

		result := chatList.CreateChat(name)
		expected := &Chat{
			ID: name,
		}
		if result.ID != expected.ID {
			t.Errorf("Fail")
		}
	})
}

func TestConcurrentCreateChat(t *testing.T) {
	wg := sync.WaitGroup{}
	var chatList ChatList
	chatList.Chats = make(map[string]*Chat)
	chatNames := []string{"cool chat", "foo", "mooze", "stalker", "4335", "grAWSRW", "GRERHAS", "h3af3"}

	for i := 0; i < len(chatNames); i++ {
		name := chatNames[i]
		wg.Add(1)
		go func() {
			chatList.CreateChat(name)
			wg.Done()
		}()
		wg.Wait()
	}

	for _, key := range chatNames {
		_, exists := chatList.Chats[key]
		if !exists {
			t.Errorf("Error during concurrent chat creation, chat '%s' must exist", key)
		}
	}
	if len(chatNames) != len(chatList.Chats) {
		t.Errorf("Error during concurrent chat creation, length of chat capacity should be %d, instead got %d", len(chatNames), len(chatList.Chats))
	}
}

type MockedDatabase struct{}

type MockedChatHandler interface {
	GetChats() ([]string, error)
}

func (mdb MockedDatabase) GetAllChats() ([]string, error) {
	chatNames := []string{"chat 1", "chat 2"}
	return chatNames, nil
}

func TestInitializeRooms(t *testing.T) {
	var chatList ChatList
	mockedDB := MockedDatabase{}
	wantKeys := []string{"chat 1", "chat 2"}

	chatList.initializeRooms(mockedDB)

	if len(chatList.Chats) != len(wantKeys) {
		t.Errorf("Expected %d chats, got %d", len(wantKeys), len(chatList.Chats))
	}

	for _, key := range wantKeys {
		if _, exists := chatList.Chats[key]; !exists {
			t.Errorf("Expected chatList to contain key %q", key)
		}
	}

}

func TestAddClientToSubRooms(t *testing.T) {
	var chatList ChatList
	chatList.Chats = make(map[string]*Chat)

	chatList.CreateChat("chat 1")
	subRooms := []string{"chat 1"}

	client := &Client{
		subs: subRooms,
	}

	chatList.addClientToSubRooms(client)
	expectedCapacity := 1
	actualCapacity := len(chatList.Chats["chat 1"].members)

	if expectedCapacity != actualCapacity {
		t.Errorf("Expected chat to have %d members, but got %d", expectedCapacity, actualCapacity)
	}
}

func TestConcurrentAddClientToSubRooms(t *testing.T) {
	var chatList ChatList
	chatList.Chats = map[string]*Chat{"chat 1": &Chat{ID: "chat 1"}}
	wg := sync.WaitGroup{}
	clients := []*Client{}
	names := []string{"joe", "bill", "grigorovich", "mitya"}
	sub := []string{"chat 1"}
	for i := 0; i < len(names); i++ {
		clients = append(clients, &Client{
			username: names[i],
			subs:     sub,
		})
	}

	for _, client := range clients {
		wg.Add(1)
		go func(client *Client) {
			defer wg.Done()
			chatList.addClientToSubRooms(client)
		}(client)
	}
	wg.Wait()

	specifiedChat := chatList.Chats["chat 1"]
	if len(specifiedChat.members) != len(names) {
		t.Errorf("Expected %d members, got %d", len(names), len(specifiedChat.members))
	}

	for _, name := range names {

		var contains = func(name string) bool {
			for _, cl := range specifiedChat.members {
				if cl.username == name {
					return true
				}
			}
			return false
		}
		ok := contains(name)

		if !ok {
			t.Errorf("Expected chat to contain member with username %s", name)
		}
	}

}
