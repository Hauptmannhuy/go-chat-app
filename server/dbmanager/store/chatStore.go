package store

import (
	"fmt"
	"go-chat-app/dbmanager/errordb"
	"log"
)

type ChatStore interface {
	LoadSubscribedChats(username string) ([]interface{}, error)
	SaveChat(name, creatorID string) (string, error)
	SavePrivateChat(u1id, u2id string) (interface{}, error)
	GetChats() (Chats, error)
	SearchChat(input, userID string) (interface{}, error)
	LoadSubscribedPrivateChats(id string) (interface{}, error)
	RetrieveGroupChatCreatorID(chatID string) string
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
	}
	return result, nil
}

func (s *SQLstore) SaveChat(name, creatorID string) (string, error) {

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

	return id, err
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
		results[data.Name] = data
	}
	return results, err
}

func (s *SQLstore) SavePrivateChat(user1id, user2id string) (interface{}, error) {
	user1name := s.retrieveUsername(user1id)
	user2name := s.retrieveUsername(user2id)
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

	err = res.Scan(&data.ChatName, &data.ChatID)
	if err != nil {
		log.Fatal("error scanning in savePrivateChat", err)
	}

	tr.Commit()
	return data, nil
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

func (s *SQLstore) RetrieveGroupChatCreatorID(chatID string) string {
	var creatorID string
	err := s.DB.QueryRow(`
		SELECT creator_id FROM group_chats WHERE id = $1
	`, chatID).Scan(&creatorID)
	if err != nil {
		fmt.Println("Error retrieving group creator ID", err)
		return ""
	}
	return creatorID
}
