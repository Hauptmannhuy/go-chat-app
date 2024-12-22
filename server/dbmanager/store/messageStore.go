package store

import (
	"fmt"
	"go-chat-app/dbmanager/errordb"
)

type MessageStore interface {
	SaveMessage(body, chatID, userID string) (int, error)
	GetChatsMessages(subs []string) (interface{}, error)
}

func (s *SQLstore) SaveMessage(body, chatName, userID string) (int, error) {
	messageID, err := s.retrieveLastMessageID(chatName)
	if err != nil {
		return -1, err
	}
	_, err = s.DB.Exec("INSERT INTO messages (body, user_id, chat_name, message_id) VALUES ($1, $2, $3, $4)", body, userID, chatName, messageID)
	if err != nil {
		fmt.Println("failed to save message")
		return -1, err
	}
	return messageID, err
}

func (s *SQLstore) GetChatsMessages(subs []string) (interface{}, error) {
	chats := make(map[string][]interface{})
	queryResults := make(map[string]bool)
	for _, key := range subs {
		queryResults[key] = false
	}

	for _, sub := range subs {

		query := `SELECT m.message_id, m.body, m.chat_name, m.user_id, u.username 
							FROM messages AS m
							JOIN users AS u
							ON m.user_id = u.id
							WHERE chat_name = $1`

		rows, err := s.DB.Query(query, sub)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var resultRows []interface{}
		for rows.Next() {
			queryResults[sub] = true
			var chatMessage Message

			if err := rows.Scan(&chatMessage.MessageID, &chatMessage.Body, &chatMessage.ChatID, &chatMessage.UserID, &chatMessage.Username); err != nil {
				return nil, err
			}
			resultRows = append(resultRows, chatMessage)

		}

		chats[sub] = resultRows
	}

	for key, keyVal := range queryResults {
		if !keyVal {
			delete(chats, key)
		}
	}

	if len(chats) == 0 {
		return nil, &errordb.ErrorDB{"No messages to fetch"}
	}

	return chats, nil
}
