package main

import "fmt"

type broadcastHandler interface {
	exec(sender *Client, payload interface{})
}

type newMsgAlgo struct{}
type newDialogueAlgo struct{}
type newGroupChatAlgo struct{}
type newSubAlgo struct{}
type searchQueryAlgo struct{}
type chatStoreAlgo struct{}
type messageStoreAlgo struct{}
type errorAlgo struct{}

func (newMsgAlgo *newMsgAlgo) exec(sender *Client, payload interface{}) {
	userMsg, _ := payload.(*UserMessage)
	destination := userMsg.ChatName
	hub := sender.hub
	envelope := OutEnvelope{
		Type: "NEW_MESSAGE",
		Data: userMsg,
	}
	room, ok := hub.roomRegister.rooms[destination]
	if ok {
		room.read <- envelope
	} else {
		fmt.Printf("room %s does not exist", destination)
	}
}

func (newMsgAlgo *newDialogueAlgo) exec(sender *Client, payload interface{}) {
	newDialog, _ := payload.(*NewPrivateChat)
	hub := sender.hub
	destination := newDialog.ChatName
	fmt.Println("destination", destination)
	hub.roomRegister.addRoom(destination)
	initiator, receiver := hub.connections[newDialog.InitiatorID], hub.connections[newDialog.ReceiverID]
	hub.roomRegister.addClient(destination, initiator)
	if _, online := hub.connections[newDialog.ReceiverID]; online {
		hub.roomRegister.addClient(destination, receiver)
	}

	newDialogEnv := OutEnvelope{
		Type: "NEW_PRIVATE_CHAT",
		Data: newDialog,
	}
	newMsgEnv := OutEnvelope{
		Type: "NEW_MESSAGE",
		Data: UserMessage{
			Body:      newDialog.Message,
			UserID:    newDialog.InitiatorID,
			ChatName:  newDialog.ChatName,
			Username:  newDialog.Username,
			MessageID: 0,
		},
	}

	hub.roomRegister.rooms[destination].read <- newDialogEnv
	hub.roomRegister.rooms[destination].read <- newMsgEnv
}

func (newMsgAlgo *newGroupChatAlgo) exec(sender *Client, payload interface{}) {
	hub := sender.hub
	data, _ := payload.(*NewGroupChat)
	destination := data.Name
	outEnv := OutEnvelope{
		Type: "NEW_GROUP_CHAT",
		Data: payload,
	}
	hub.roomRegister.addRoom(destination)
	hub.roomRegister.addClient(destination, sender)
	sender.messageBuffer <- outEnv
}

func (newMsgAlgo *newSubAlgo) exec(sender *Client, payload interface{}) {
	hub := sender.hub
	data, _ := payload.(*Subscription)
	destination := data.ChatName
	chatEnv := OutEnvelope{
		Type: "JOIN_CHAT",
		Data: payload,
	}
	msgEnv := OutEnvelope{
		Type: "NEW_MESSAGE",
		Data: UserMessage{
			Body:      data.BodyMessage,
			ChatName:  data.ChatName,
			Username:  data.Username,
			MessageID: data.msgID,
		},
	}
	_, ok := hub.roomRegister.rooms[destination]
	if !ok {
		hub.roomRegister.addRoom(destination)
	}
	hub.roomRegister.addClient(destination, sender)
	hub.roomRegister.rooms[destination].read <- msgEnv
	sender.messageBuffer <- chatEnv
}

func (newMsgAlgo *searchQueryAlgo) exec(sender *Client, payload interface{}) {
	outEnv := OutEnvelope{
		Type: "SEARCH_QUERY",
		Data: payload,
	}
	sender.messageBuffer <- outEnv
}

func (newMsgAlgo *chatStoreAlgo) exec(sender *Client, payload interface{}) {
	data := payload.(*WebSocketChatStore)

	outEnv := OutEnvelope{
		Type: "LOAD_SUBS",
		Data: data.Data,
	}
	sender.messageBuffer <- outEnv
}

func (newMsgAlgo *messageStoreAlgo) exec(sender *Client, payload interface{}) {
	data := payload.(*WebSocketMessageStore)

	outEnv := OutEnvelope{
		Type: "LOAD_MESSAGES",
		Data: data.Data,
	}
	sender.messageBuffer <- outEnv
}

func (errorAlgo *errorAlgo) exec(sender *Client, payload interface{}) {

}
