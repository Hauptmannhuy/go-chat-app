package main

import (
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
