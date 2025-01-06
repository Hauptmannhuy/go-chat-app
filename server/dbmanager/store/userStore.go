package store

import (
	"fmt"
	"go-chat-app/dbmanager/errordb"
)

type UserStore interface {
	SaveAccount(name, email, pass string) (string, error)
	AuthenticateAccount(name, pass string) (string, error)
	SearchUser(username, userID string) (map[string]UserContainerData, error)
}

func (s *SQLstore) SearchUser(input, userID string) (map[string]UserContainerData, error) {
	userName := s.retrieveUsername(userID)
	userMap := make(map[string]UserContainerData)
	query := `
	SELECT u.username, u.id, 
	CASE	
		WHEN pc.id IS NULL THEN -1
		ELSE pc.id
	END AS private_chat_id
	, CASE 
		WHEN pc.id IS NULL THEN ''
		ELSE pc.chat_name
	END AS chat_name
	, CASE 
		WHEN pc.id IS NULL THEN -1
		ELSE pc.user1_id
	END AS user1_id
	, CASE
		WHEN pc.user2_id IS NULL THEN -1
		ELSE pc.user2_id
	END AS user2_id
	,
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

		resultChat := PrivateChatInfo{
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
		userData := UserContainerData{}

		userData.Profile = resultProfile
		userData.Chat = resultChat
		userMap[resultProfile.Username] = userData

	}
	return userMap, nil

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
