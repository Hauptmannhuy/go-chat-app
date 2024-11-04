package store

import (
	"database/sql"
	"fmt"
	"go-chat-app/dbmanager/errordb"
	"log"

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
	ID       string `json:"id"`
}

type loginUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

type SearchResults []interface{}

type UserStore interface {
	SaveAccount(name, email, pass string) (string, error)
	AuthenticateAccount(name, pass string) (string, error)
	SearchUser(username string) (interface{}, error)
}

type MessageStore interface {
	SaveMessage(body, chatID, userID string) error
	GetChatsMessages(subs []string) (interface{}, error)
}

type ChatStore interface {
	LoadSubscribedChats(username string) ([]interface{}, error)
	SaveChat(id, chat_type string) error
	SavePrivateChat(u1id, u2id string) (string, error)
	GetChats() ([]string, error)
	SearchChat(input, userID string) ([]interface{}, error)
}

type SubscriptionStore interface {
	LoadSubscriptions(username string) ([]string, error)
	SaveSubscription(username, chatID string) error
}

type SQLstore struct {
	DB *sql.DB
}

func (s *SQLstore) GetChats() ([]string, error) {
	rows, err := s.DB.Query(`SELECT chat_id FROM chats`)
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

func (s *SQLstore) SaveChat(id, chat_type string) error {
	tr, _ := s.DB.Begin()

	_, err := s.DB.Exec(`
		INSERT INTO chats (chat_id, chat_type)
			VALUES ($1, $2)
	`, id, chat_type)

	if err != nil {
		fmt.Println(err)
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
		return nil, &errordb.ErrorDB{"No messages to fetch"}
	}

	return chats, nil
}

func (s *SQLstore) AuthenticateAccount(name, pass string) (string, error) {

	var data loginUserData
	row := s.DB.QueryRow(
		`
	SELECT id, username, password FROM users
	WHERE username = $1
	`, name)

	fmt.Println(row)
	err := row.Scan(&data.ID, &data.Username, &data.Password)
	if err != nil {
		fmt.Println(err)
		err = errordb.ParseError(err.Error())
		return "", err
	}

	passValid, err := authenticatePass([]byte(data.Password), []byte(pass))

	if passValid {
		return data.ID, nil
	} else {
		err = errordb.ParseError(err.Error())
		return "", err
	}
}

func (s *SQLstore) SaveAccount(name, email, pass string) (string, error) {
	var id string
	encryptedPass, err := encryptPassword(pass)
	if err != nil {
		fmt.Println((err))
		return "", err
	}

	tr, err := s.DB.Begin()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	res, err := tr.Exec(`
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3)
	`, name, email, encryptedPass)
	fmt.Println("result reg:", res)

	if err != nil {
		fmt.Println(err)
		err = errordb.ParseError(err.Error())
		tr.Rollback()
		return "", err
	}
	tr.Commit()

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	row := s.DB.QueryRow(`
		SELECT last_value from users_id_seq
	`)
	row.Scan(&id)
	fmt.Println("LAST REGISTERED ID:", id)
	return id, nil
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

func (s *SQLstore) SearchChat(input, userID string) ([]interface{}, error) {

	var results []interface{}
	query := fmt.Sprintf(`
		SELECT chat_id, chat_type
		FROM chats
		WHERE chat_id NOT IN (SELECT chat_id FROM subscriptions WHERE username = '%s')
		AND chat_type = 'group' AND chat_id ILIKE '%s'
		LIMIT 25
	`, userID, input+"%")
	fmt.Println(query)
	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Println("Error during search query:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data struct {
			Type string `json:"chat_type"`
			ID   string `json:"chat_id"`
		}
		rows.Scan(&data.ID, &data.Type)
		fmt.Println(data)
		results = append(results, data)
	}
	return results, err
}

func (s *SQLstore) SearchUser(input string) (interface{}, error) {
	userMap := make(map[string]interface{})
	query := fmt.Sprintf(`SELECT username, id FROM users WHERE username ILIKE '%s' LIMIT 10`, input+"%")
	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result UserProfileData
		fmt.Println(result)
		rows.Scan(&result.Username, &result.ID)
		userMap[result.Username] = result
	}
	return userMap, nil

}

func (s *SQLstore) LoadSubscribedChats(username string) ([]interface{}, error) {

	var res SearchResults

	tr, _ := s.DB.Begin()

	rows, err := tr.Query(`
	SELECT c.chat_id, c.chat_type FROM chats AS c
	JOIN subscriptions AS s
	ON s.chat_id = c.chat_id
	WHERE s.username = $1
	`, username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var data struct {
			ChatID   string `json:"chat_id"`
			ChatType string `json:"chat_type"`
		}
		rows.Scan(&data.ChatID, &data.ChatType)
		res = append(res, data)
	}
	fmt.Println("Subs:", res)
	return res, nil

}

func (s *SQLstore) SavePrivateChat(user1id, user2id string) (string, error) {
	user1name := s.retrieveUsername(user1id)
	user2name := s.retrieveUsername(user2id)
	fmt.Println("input:", user1id, user2id, user1name, user2name)
	chatID := user1name + "_" + user2name
	tr, _ := s.DB.Begin()
	query := `
		INSERT INTO private_chats (user1_id, user2_id, chat_id) 
		VALUES ($1, $2, $3)
	`
	_, err := tr.Exec(query, user1id, user2id, chatID)
	if err != nil {
		fmt.Println(err)
		tr.Rollback()
		return "", err
	}
	tr.Commit()
	return chatID, nil
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
