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
	Username  string `json:"username"`
	Body      string `json:"body"`
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_name"`
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

type privateChat struct {
	ChatName  string `json:"chat_name"`
	ChatID    int    `json:"chat_id"`
	User1ID   int    `json:"user1_id"`
	User2ID   int    `json:"user2_id"`
	Handshake bool   `json:"handshake"`
}

type PrivateChatInfo struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ChatInfo struct {
	Name        string
	ID          int
	Subscribers []string
	ChatType    string
}

type Chats map[string]ChatInfo

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
	SavePrivateChat(u1id, u2id string) (interface{}, error)
	GetChats() (Chats, error)
	SearchChat(input, userID string) (interface{}, error)
	LoadSubscribedPrivateChats(id string) (interface{}, error)
}

type SubscriptionStore interface {
	LoadSubscriptions(username string) ([]string, error)
	SaveSubscription(userID, chatID string) error
}

type SQLstore struct {
	DB *sql.DB
}

func (s *SQLstore) GetChats() (Chats, error) {
	var result = make(Chats)

	combinedRows, err := s.DB.Query(
		`SELECT chat_name, id, 'private' AS chat_type
		FROM private_chats
		UNION ALL
		SELECT chat_name, id, 'group' AS chat_type
		FROM group_chats;
	`)

	if err != nil {
		fmt.Println("Error getting private chats", err)
		return nil, err
	}

	for combinedRows.Next() {
		fmt.Println("got")
		var name string
		var chatType string
		var ID int

		err := combinedRows.Scan(&name, &ID, &chatType)
		if err != nil {
			fmt.Println("error during scan of getting chats", err)
		}
		var subscribers []string
		if chatType == "private" {
			rows, err := s.DB.Query(`
			SELECT u1.username AS username1, u2.username AS username2
			FROM private_chats AS prc
			JOIN users AS u1
			ON u1.id = prc.user1_id
			JOIN users AS u2
			ON u2.id = prc.user2_id
			WHERE prc.id = $1 `, ID)

			if err != nil {
				log.Fatal("error in GetChats:", err)
			}

			defer rows.Close()
			for rows.Next() {
				var username1 string
				var username2 string
				err := rows.Scan(&username1, &username2)
				if err != nil {
					log.Fatal("error in GetChats:", err)
				}
				subscribers = append(subscribers, username1, username2)

			}
			result[name] = ChatInfo{
				Name:        name,
				ID:          ID,
				ChatType:    chatType,
				Subscribers: subscribers,
			}

		} else {
			rows, err := s.DB.Query(`
			SELECT u.username 
			FROM users AS u
			JOIN group_chat_subs AS gcs
			ON gcs.user_id = u.id
			WHERE gcs.chat_id = $1 `, ID)

			if err != nil {
				log.Fatal("error in GetChats:", err)
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				err := rows.Scan(&name)
				if err != nil {
					log.Fatal("error in GetChats:", err)
				}
				subscribers = append(subscribers, name)
			}
		}
		result[name] = ChatInfo{
			Name:        name,
			ID:          ID,
			ChatType:    chatType,
			Subscribers: subscribers,
		}
		fmt.Println(result[name])
	}
	return result, nil
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

func (s *SQLstore) SaveMessage(body, chatName, userID string) error {
	messageID, err := s.retrieveLastMessageID(chatName)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec("INSERT INTO messages (body, user_id, chat_name, message_id) VALUES ($1, $2, $3, $4)", body, userID, chatName, messageID)
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

	rows, err := s.DB.Query(`
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
		var subGroup string
		rows.Scan(&subGroup)
		res = append(res, subGroup)
	}

	rows, err = s.DB.Query(`
	SELECT chat_name
	FROM private_chats
	WHERE user1_id = $1 OR user2_id = $1
	`, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var subPrivate string
		rows.Scan(&subPrivate)
		res = append(res, subPrivate)
	}

	fmt.Println("Subs:", res)
	return res, nil

}

func (s *SQLstore) SaveSubscription(userID, chatID string) error {
	fmt.Println("393", userID, chatID)
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

func (s *SQLstore) SearchChat(input, userID string) (interface{}, error) {

	results := make(map[string]interface{})
	query := fmt.Sprintf(`
		SELECT gc.chat_name, gc.id,
		CASE 
			WHEN gcs.id IS NOT NULL THEN TRUE
			ELSE FALSE
		END AS is_subscribed
		FROM group_chats AS gc
		LEFT JOIN group_chat_subs AS gcs
		ON gcs.user_id = %s AND gc.id = gcs.chat_id
		WHERE gc.chat_name ILIKE '%s'
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
			Name         string `json:"chat_name"`
			ID           string `json:"chat_id"`
			IsSubscribed string `json:"is_subscribed"`
		}
		rows.Scan(&data.Name, &data.ID, &data.IsSubscribed)
		fmt.Println(data)
		results[data.Name] = data
	}
	return results, err
}

func (s *SQLstore) SearchUser(input, userID string) (interface{}, error) {
	userName := s.retrieveUsername(userID)
	userMap := make(map[string]interface{})
	query := `
	SELECT u.username, u.id, pc.id, pc.chat_name, pc.user1_id, pc.user2_id,
	CASE
		WHEN pc.id IS NOT NULL THEN TRUE
		ELSE FALSE
		END AS handshake
	FROM users AS u
	LEFT JOIN private_chats AS pc
	ON (CAST($1 AS INTEGER) = pc.user1_id AND u.id = pc.user2_id)
	 OR (u.id = pc.user1_id AND CAST($1 AS INTEGER) = pc.user2_id)
	WHERE u.username ILIKE $2 AND u.username NOT ILIKE $3 LIMIT 10`

	rows, err := s.DB.Query(query, userID, input+"%", userName)

	fmt.Println(rows)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var username, chatName string
		var id, pcID, user1ID, user2ID int
		var handshake bool

		if err := rows.Scan(&username, &id, &pcID, &chatName, &user1ID, &user2ID, &handshake); err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}

		fmt.Printf("username: %s, id: %d, pcID: %d, chatName: %s, user1ID: %d, user2ID: %d, handshake: %v\n",
			username, id, pcID, chatName, user1ID, user2ID, handshake)

		resultChat := privateChat{
			ChatName:  chatName,
			ChatID:    pcID,
			User1ID:   user1ID,
			User2ID:   user2ID,
			Handshake: handshake,
		}

		resultProfile := UserProfileData{
			Username: username,
			ID:       id,
		}
		var container struct {
			Profile UserProfileData `json:"profile"`
			Chat    privateChat     `json:"chat"`
		}

		container.Profile = resultProfile
		container.Chat = resultChat
		userMap[resultProfile.Username] = container
		fmt.Println(container)
	}
	return userMap, nil

}

func (s *SQLstore) LoadSubscribedChats(userID string) ([]interface{}, error) {

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

func (s *SQLstore) SavePrivateChat(user1id, user2id string) (interface{}, error) {
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
		return nil, err
	}

	res := tr.QueryRow(`
		SELECT chat_name, id 
		FROM private_chats
		WHERE chat_name = $1
	`, chatName)

	var data PrivateChatInfo

	err = res.Scan(&data.Name, &data.ID)
	if err != nil {
		log.Fatal("error scanning in savePrivateChat", err)
	}

	tr.Commit()
	return data, nil
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
