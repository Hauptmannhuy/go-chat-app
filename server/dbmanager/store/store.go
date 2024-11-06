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

type privateChat struct {
	ChatName string `json:"chat_name"`
	ChatID   string `json:"chat_id"`
	User1ID  string `json:"user1_id"`
	User2ID  string `json:"user2_id"`
}

type SearchResults []interface{}

type UserStore interface {
	SaveAccount(name, email, pass string) (string, error)
	AuthenticateAccount(name, pass string) (string, error)
	SearchUser(username, userID string) (interface{}, error)
}

type MessageStore interface {
	SaveMessage(body, chatID, userID string) error
	GetChatsMessages(subs []string) (interface{}, error)
}

type ChatStore interface {
	LoadSubscribedChats(username string) ([]interface{}, error)
	SaveChat(name, creatorID string) (string, error)
	SavePrivateChat(u1id, u2id string) (string, error)
	GetChats() ([]string, error)
	SearchChat(input, userID string) ([]interface{}, error)
	LoadSubscribedPrivateChats(id string) (interface{}, error)
}

type SubscriptionStore interface {
	LoadSubscriptions(username string) ([]string, error)
	SaveSubscription(userID, chatID string) error
}

type SQLstore struct {
	DB *sql.DB
}

func (s *SQLstore) GetChats() ([]string, error) {
	combinedRows, err := s.DB.Query(
		`SELECT chat_name, 'private' AS chat_type
		FROM private_chats
		UNION ALL
		SELECT chat_name, 'group' AS chat_type
		FROM group_chats;
	`)
	if err != nil {
		fmt.Println("Error getting private chats", err)
		return nil, err
	}

	var res []string

	for combinedRows.Next() {
		fmt.Println("got")
		var name string
		var chatType string
		err := combinedRows.Scan(&name, &chatType)
		if err != nil {
			fmt.Println("error during scan of getting chats", err)
		}
		res = append(res, name)
	}
	fmt.Println("res", res)
	return res, nil
}

func (s *SQLstore) SaveChat(name, creatorID string) (string, error) {

	fmt.Println(name, creatorID)
	tr, _ := s.DB.Begin()

	_, err := s.DB.Exec(`
		INSERT INTO group_chats (chat_name, creator_id)
			VALUES ($1, $2)
	`, name, creatorID)

	if err != nil {
		fmt.Println("error saving group chat", err)
		err = errordb.ParseError(err.Error())
		tr.Rollback()
		return "", err
	}
	tr.Commit()
	id := s.retrieveGroupChatIndex(name)
	fmt.Println("group id", id)

	return id, err
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
		fmt.Println(err)
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
		fmt.Println("error saving account", err)
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

func (s *SQLstore) LoadSubscriptions(userID string) ([]string, error) {

	var res []string

	tr, _ := s.DB.Begin()

	rows, err := tr.Query(`
	SELECT gc.chat_name 
	FROM group_chats AS gc
	JOIN group_chat_subs AS gcs
	ON gc.id = gcs.chat_id
	WHERE gcs.user_id = $1

	`, userID)
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

func (s *SQLstore) SaveSubscription(userID, chatID string) error {
	tr, _ := s.DB.Begin()
	_, err := tr.Exec(`
		INSERT INTO group_chat_subs (user_id, chat_id) VALUES ($1, $2)
	`, userID, chatID)

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
		SELECT chat_name, id
		FROM group_chats
		WHERE chat_name NOT IN (SELECT chat_name FROM group_chat_subs WHERE user_id = %s)
	  AND chat_name ILIKE '%s'
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
			Name string `json:"chat_name"`
			ID   string `json:"chat_id"`
		}
		rows.Scan(&data.Name, &data.ID)
		fmt.Println(data)
		results = append(results, data)
	}
	return results, err
}

func (s *SQLstore) SearchUser(input, userID string) (interface{}, error) {
	userID = s.retrieveUsername(userID)
	userMap := make(map[string]interface{})
	query := fmt.Sprintf(`SELECT u.username, u.id, pc.id, pc.chat_name, pc.user1_id, pc.user2_id  FROM users AS u
	LEFT JOIN private_chats AS pc
	ON u.id = pc.user1_id OR u.id = pc.user2_id
	WHERE u.username ILIKE '%s' AND u.username NOT ILIKE '%s' LIMIT 10`, input+"%", userID)

	fmt.Println(query)

	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resultChat privateChat
		var resultProfile UserProfileData
		var container struct {
			Profile UserProfileData `json:"profile"`
			Chat    privateChat     `json:"chat"`
		}

		rows.Scan(&resultProfile.Username, &resultProfile.ID, &resultChat.ChatID, &resultChat.ChatName, &resultChat.User1ID, &resultChat.User2ID)

		container.Profile = resultProfile
		userMap[resultProfile.Username] = container
		fmt.Println(container)
	}
	return userMap, nil

}

func (s *SQLstore) LoadSubscribedChats(userID string) ([]interface{}, error) {

	fmt.Println("id", userID)
	var res SearchResults
	tr, _ := s.DB.Begin()
	rows, err := tr.Query(`
	SELECT c.id, c.chat_name FROM group_chats AS c
	JOIN group_chat_subs AS s
	ON s.chat_id = c.id
	WHERE s.user_id = $1
	`, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var data struct {
			ChatID   string `json:"chat_id"`
			ChatName string `json:"chat_name"`
		}
		err := rows.Scan(&data.ChatID, &data.ChatName)
		if err != nil {
			log.Fatal("Error scanning subscribed group chats", err)
		}
		res = append(res, data)
	}
	return res, nil

}

func (s *SQLstore) LoadSubscribedPrivateChats(userID string) (interface{}, error) {
	privateChatsMap := make(map[string]interface{})
	query := fmt.Sprintf(`
		SELECT id, user1_id, user2_id, chat_name FROM private_chats WHERE user1_id = '%s' OR user2_id = '%s'
	`, userID, userID)
	rows, err := s.DB.Query(query)

	if err != nil {
		fmt.Println("Error loading private chats", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var resultRow struct {
			ChatID   string `json:"chat_id"`
			ChatName string `json:"chat_name"`
			User1ID  string `json:"user1_id"`
			User2ID  string `json:"user2_id"`
		}
		err := rows.Scan(&resultRow.ChatID, &resultRow.User1ID, &resultRow.User2ID, &resultRow.ChatName)
		if err != nil {
			fmt.Println("Error loading private chats", err)
			return nil, err
		}
		privateChatsMap[resultRow.ChatID] = resultRow
	}
	return privateChatsMap, nil
}

func (s *SQLstore) SavePrivateChat(user1id, user2id string) (string, error) {
	user1name := s.retrieveUsername(user1id)
	user2name := s.retrieveUsername(user2id)
	fmt.Println("input:", user1id, user2id, user1name, user2name)
	chatName := user1name + "_" + user2name
	tr, _ := s.DB.Begin()
	query := `
		INSERT INTO private_chats (user1_id, user2_id, chat_name) 
		VALUES ($1, $2, $3)
	`
	_, err := tr.Exec(query, user1id, user2id, chatName)
	if err != nil {
		fmt.Println("Error saving private chats", err)
		tr.Rollback()
		return "", err
	}
	tr.Commit()
	return chatName, nil
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
