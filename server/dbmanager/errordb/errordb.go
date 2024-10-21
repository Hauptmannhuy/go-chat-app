package errordb

import (
	"fmt"
)

type ErrorDB struct {
	Message string
}

func (e *ErrorDB) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func ParseError(dberr string) error {
	var keyErrors = map[string]string{
		`pq: duplicate key value violates unique constraint "users_username_key"`: "This name is occupied",
		`pq: duplicate key value violates unique constraint "users_email_key"`:    "This email is already occupied",
		`pq: duplicate key value violates unique constraint "chats_pkey"`:         "This chat name is already occupied",
		"sql: no rows in result set":                                              "Account with provided data doesn't exists.",
		`crypto/bcrypt: hashedPassword is not the hash of the given password`:     "Wrong password.",
	}

	el, ok := keyErrors[dberr]

	if ok {
		return &ErrorDB{el}
	} else {
		return &ErrorDB{"Unknown database error"}
	}
}
