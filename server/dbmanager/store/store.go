package store

import (
	"database/sql"
	"fmt"
	"go-chat-app/dbmanager/errordb"

	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	ChatID string `json:"chatID"`
	// CreatedAt string `json:"created_at"`
}

type loginUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	SaveAccount(name, email, pass string) error
	AuthenticateAccount(name, pass string) error
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

func (s *SQLstore) retrieveLastMessageID(chatID string) (int, error) {
	tr, _ := s.DB.Begin()

	_, err := s.DB.Exec(`
		INSERT INTO last_messages_ids (chatID, last_message_id)
        VALUES ($1, 0)
        ON CONFLICT (chatID) DO UPDATE SET last_message_id = last_messages_ids.last_message_id + 1
		`, chatID)
	if err != nil {
		fmt.Println(err)
		tr.Rollback()
		return 0, err
	}

	var message_id int
	err = s.DB.QueryRow(`
		SELECT last_message_id FROM last_messages_ids WHERE chatID = $1
	`, chatID).Scan(&message_id)
	return message_id, err
}

func (s *SQLstore) SaveMessage(body, chatID string) error {
	messageID, err := s.retrieveLastMessageID(chatID)

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

func (s *SQLstore) AuthenticateAccount(name, pass string) error {

	tr, _ := s.DB.Begin()

	var data loginUserData
	row := tr.QueryRow(
		`
	SELECT username, password FROM users
	WHERE username = $1
	`, name)

	fmt.Println(row)
	err := row.Scan(&data.Username, &data.Password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(data)
	passValid, err := authenticatePass([]byte(data.Password), []byte(pass))

	if passValid {
		tr.Commit()
		return nil
	} else {
		return err
	}
}

func (s *SQLstore) SaveAccount(name, email, pass string) error {

	encryptedPass, err := encryptPassword(pass)
	if err != nil {
		fmt.Println((err))
		return err
	}

	tr, err := s.DB.Begin()

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = tr.Exec(`
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3)
	`, name, email, encryptedPass)

	if err != nil {
		fmt.Println(err)
		err = errordb.ParseError(err.Error())
		tr.Rollback()
		return err
	}
	tr.Commit()
	return nil
}

func encryptPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(hash), err
}

func authenticatePass(hashedPass, pass []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPass, pass)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
