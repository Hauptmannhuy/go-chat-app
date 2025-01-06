package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Body      string `json:"body"`
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_name"`
}

type UserContainerData struct {
	Profile UserProfileData `json:"profile"`
	Chat    PrivateChatInfo `json:"chat"`
}

type UserProfileData struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}

type loginUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

// type PrivateChat struct {
// 	ChatName  string `json:"chat_name"`
// 	ChatID    int    `json:"chat_id"`
// 	User1ID   int    `json:"user1_id"`
// 	User2ID   int    `json:"user2_id"`
// 	Handshake bool   `json:"handshake"`
// }

type PrivateChatInfo struct {
	ChatName  string `json:"chat_name"`
	ChatID    int    `json:"chat_id"`
	User1ID   int    `json:"user1_id"`
	User2ID   int    `json:"user2_id"`
	Handshake bool   `json:"handshake"`
}

type ChatInfo struct {
	Name        string
	ID          int
	Subscribers []string
	ChatType    string
}

type Chats map[string]ChatInfo

type SearchResults []interface{}

type SQLstore struct {
	DB *sql.DB
}

func (s *SQLstore) retrieveLastMessageID(chatID string) (int, error) {
	tr, _ := s.DB.Begin()

	_, err := s.DB.Exec(`
		INSERT INTO last_messages_ids (chat_id, last_message_id)
        VALUES ($1, 0)
        ON CONFLICT (chat_id) DO UPDATE SET last_message_id = last_messages_ids.last_message_id + 1
		`, chatID)
	if err != nil {
		fmt.Println(err)
		tr.Rollback()
		return 0, err
	}

	var message_id int
	err = s.DB.QueryRow(`
		SELECT last_message_id FROM last_messages_ids WHERE chat_id = $1
	`, chatID).Scan(&message_id)
	return message_id, err
}

func encryptPassword(pass string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	fmt.Println("registered pass", hash)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return hash, err
}

func authenticatePass(hashedPass, pass []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPass, pass)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func (s *SQLstore) retrieveUsername(id string) string {
	var username string
	query := `
		SELECT username FROM users
		WHERE $1 = id
	`
	row := s.DB.QueryRow(query, id)
	err := row.Scan(&username)
	if err != nil {
		fmt.Println(err)
	}
	return username
}

func (s *SQLstore) retrieveGroupChatIndex(name string) string {
	var id string
	query := `
		SELECT id FROM group_chats
		WHERE chat_name = $1
	`
	row := s.DB.QueryRow(query, name)
	err := row.Scan(&id)
	if err != nil {
		log.Fatal("Error retrieving group chat column index:", err)
	}
	return id
}
