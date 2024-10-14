package store

import (
	"database/sql"
	"fmt"
	"go-chat-app/dbmanager/errordb"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Message struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	ChatID string `json:"chatID"`
	// CreatedAt string `json:"created_at"`
}

type UserStore interface {
	SaveAccount(name, email, pass string) error
}

type MessageStore interface {
	SaveMessage(body, description string) error
	GetAllMessages() ([]Message, error)
}

type ChatStore interface {
	SaveChat(ID string) error
}

type SQLstore struct {
	DB *sql.DB
}

func (s *SQLstore) SaveChat(ID string) error {
	tr, _ := s.DB.Begin()

	_, err := s.DB.Exec(`
		INSERT INTO chats (ID)
			VALUES ($1)
	`, ID)

	if err != nil {
		err = errordb.ParseError(err.Error())
		tr.Rollback()
		return err
	}
	tr.Commit()
	return err
}

func (s *SQLstore) retrieveLastMessageID(chatID string) (error, int) {
	tr, _ := s.DB.Begin()

	_, err := s.DB.Exec(`
		INSERT INTO last_messages_ids (chatID, last_message_id)
        VALUES ($1, 0)
        ON CONFLICT (chatID) DO UPDATE SET last_message_id = last_messages_ids.last_message_id + 1
		`, chatID)
	if err != nil {
		fmt.Println(err)
		tr.Rollback()
		return err, 0
	}

	var message_id int
	err = s.DB.QueryRow(`
		SELECT last_message_id FROM last_messages_ids WHERE chatID = $1
	`, chatID).Scan(&message_id)
	return err, message_id
}

func (s *SQLstore) SaveMessage(body, chatID string) error {
	err, messageID := s.retrieveLastMessageID(chatID)

	if err != nil {
		return err
	}
	_, err = s.DB.Exec("INSERT INTO messages (body, chatID, message_id) VALUES ($1, $2, $3)", body, chatID, messageID)
	if err != nil {
		fmt.Println("failed to save message")
		return err
	}
	return err
}

func (s *SQLstore) GetAllMessages() ([]Message, error) {
	rows, err := s.DB.Query("SELECT id, body, chatID, created_at FROM messages")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var message Message

		if err := rows.Scan(&message.ID, &message.Body, &message.ChatID); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (s *SQLstore) SaveAccount(name, email, pass string) error {
	tr, err := s.DB.Begin()

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tr.Exec(`
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3)
	`, name, email, pass)

	if err != nil {
		fmt.Println(err)
		err = errordb.ParseError(err.Error())
		tr.Rollback()
		return err
	}
	tr.Commit()
	return nil
}
