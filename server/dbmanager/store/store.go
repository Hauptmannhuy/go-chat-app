package store

import (
	"database/sql"
	"fmt"
	"go-chat-app/dbmanager/errordb"

	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	UserID    string `json:"user_id"`
	Body      string `json:"body"`
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
	// CreatedAt string `json:"created_at"`
}

type UserProfileData struct {
	Username string `json:"username"`
}

type loginUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	SaveAccount(name, email, pass string) error
	AuthenticateAccount(name, pass string) error
	SearchUser(username string) (interface{}, error)
}

type MessageStore interface {
	SaveMessage(body, chatID, userID string) error
	GetChatsMessages(subs []string) (interface{}, error)
}

type ChatStore interface {
	SaveChat(ID string) error
	GetChats() ([]string, error)
	SearchChat(input string) ([]string, error)
}

type SubscriptionStore interface {
	LoadSubscriptions(username string) ([]string, error)
	SaveSubscription(username, chatID string) error
}

type SQLstore struct {
	DB *sql.DB
}

func (s *SQLstore) GetChats() ([]string, error) {
	rows, err := s.DB.Query(`SELECT id FROM chats`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var res []string
	defer rows.Close()

	for rows.Next() {
		var row string
		rows.Scan(&row)
		res = append(res, row)
	}
	return res, nil
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

func (s *SQLstore) SaveMessage(body, chatID, userID string) error {
	messageID, err := s.retrieveLastMessageID(chatID)

	if err != nil {
		return err
	}
	_, err = s.DB.Exec("INSERT INTO messages (body, user_id, chat_id, message_id) VALUES ($1, $2, $3, $4)", body, userID, chatID, messageID)
	if err != nil {
		fmt.Println("failed to save message")
		return err
	}
	return err
}

func (s *SQLstore) GetChatsMessages(subs []string) (interface{}, error) {
	chats := make(map[string][]interface{})
	queryResults := make(map[string]bool)

	for _, key := range subs {
		queryResults[key] = false
	}

	for _, sub := range subs {

		rows, err := s.DB.Query(fmt.Sprintf("SELECT message_id, body, chat_id, user_id FROM messages WHERE chat_id = '%s'", sub))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var resultRows []interface{}
		for rows.Next() {
			queryResults[sub] = true
			var chatMessage Message

			if err := rows.Scan(&chatMessage.MessageID, &chatMessage.Body, &chatMessage.ChatID, &chatMessage.UserID); err != nil {
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
		return nil, &errordb.ErrorDB{"No chats to fetch"}
	}

	return chats, nil
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
		err = errordb.ParseError(err.Error())
		return err
	}
	fmt.Println(data)
	passValid, err := authenticatePass([]byte(data.Password), []byte(pass))

	if passValid {
		tr.Commit()
		return nil
	} else {
		err = errordb.ParseError(err.Error())
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

func (s *SQLstore) LoadSubscriptions(username string) ([]string, error) {

	var res []string

	tr, _ := s.DB.Begin()

	rows, err := tr.Query(`
	SELECT chat_id from subscriptions WHERE username = $1
	`, username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var sub string
		rows.Scan(&sub)
		res = append(res, sub)
	}
	fmt.Println("Subs:", res)
	return res, nil

}

func (s *SQLstore) SaveSubscription(username, chatID string) error {
	tr, _ := s.DB.Begin()
	_, err := tr.Exec(`
		INSERT INTO subscriptions (username, chat_id) VALUES ($1, $2)
	`, username, chatID)

	if err != nil {
		fmt.Println(err)
		tr.Rollback()
		return err
	}
	tr.Commit()
	return nil
}

func (s *SQLstore) SearchChat(input string) ([]string, error) {
	var results []string
	rows, err := s.DB.Query(fmt.Sprintf(`SELECT id from chats WHERE id ILIKE '%s' LIMIT 10`, input+"%"))
	if err != nil {
		fmt.Println("Error during search query:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var res string
		rows.Scan(&res)
		results = append(results, res)
	}
	return results, err
}

func (s *SQLstore) SearchUser(input string) (interface{}, error) {
	userMap := make(map[string]interface{})
	query := fmt.Sprintf(`SELECT username FROM users WHERE username ILIKE '%s' LIMIT 10`, input+"%")
	fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result UserProfileData
		fmt.Println(result)
		rows.Scan(&result.Username)
		userMap[result.Username] = result
	}
	return userMap, nil

}
