package store

import (
	"fmt"
	"log"
)

type SubscriptionStore interface {
	LoadSubscriptions(userID int) ([]string, error)
	SaveSubscription(userID, chatID int) error
	GetPrivateChatSubs(chatName string) []int
	GetGroupChatSubs(chatName string) []int
}

func (s *SQLstore) LoadSubscriptions(userID int) ([]string, error) {

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

func (s *SQLstore) SaveSubscription(userID, chatID int) error {
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

func (s *SQLstore) GetPrivateChatSubs(chatName string) []int {
	query := `
	SELECT user1_id, user2_id
	FROM private_chats
	WHERE chat_name = $1
	`
	usersID := make([]int, 2)
	row := s.DB.QueryRow(query, chatName)
	err := row.Scan(&usersID[0], &usersID[1])

	if err != nil {
		log.Fatal("err1", err)
	}
	return usersID
}

func (s *SQLstore) GetGroupChatSubs(chatName string) []int {
	var receivers []int
	query := `
	SELECT gcs.user_id FROM group_chats AS gc
	JOIN group_chat_subs AS gcs
	ON gc.id = gcs.id
	WHERE gc.chat_name = $1`
	rows, err := s.DB.Query(query, chatName)
	if err != nil {
		log.Fatal("Error during search for subs", err)
	}
	for rows.Next() {
		var receiver int
		err := rows.Scan(&receiver)
		if err != nil {
			log.Fatal("Error during scan", err)
		}
		receivers = append(receivers, receiver)
	}
	return receivers
}
