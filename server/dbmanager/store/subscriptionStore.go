package store

import (
	"fmt"
	"log"
)

type SubscriptionStore interface {
	LoadSubscriptions(username string) ([]string, error)
	SaveSubscription(userID, chatID string) error
	GetPrivateChatSubs(chatName, sender string) []string
	GetGroupChatSubs(chatName, sender string) []string
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

func (s *SQLstore) GetPrivateChatSubs(chatName, sender string) []string {
	var receiver string
	query := `
		SELECT u.username
		FROM private_chats AS pc
		JOIN users AS u
		ON pc.user1_id = u.id OR pc.user2_id = u.id
		WHERE pc.chat_name = $1 AND u.username != $2
	`
	row := s.DB.QueryRow(query, chatName, sender)
	err := row.Scan(&receiver)

	if err != nil {
		log.Fatal(err)
	}
	return []string{receiver}
}

func (s *SQLstore) GetGroupChatSubs(chatName, sender string) []string {
	var receivers []string
	query := `
	SELECT u.username FROM users AS u
	JOIN group_chat_subs AS gcs
	ON gcs.user_id = u.id
	JOIN group_chats AS gc
	ON gc.id = gcs.chat_id
	WHERE gc.chat_name == $1 AND u.username != $2;`
	rows, err := s.DB.Query(query, chatName, sender)
	if err != nil {
		log.Fatal("Error during search for subs", err)
	}
	for rows.Next() {
		var receiver string
		err := rows.Scan(&receiver)
		if err != nil {
			log.Fatal("Error during scan", err)
		}
		receivers = append(receivers, receiver)
	}
	return receivers
}
