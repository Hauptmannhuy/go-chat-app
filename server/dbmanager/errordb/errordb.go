package errordb

import (
	"fmt"
	"strings"
)

type ErrorDB struct {
	Message string
}

func (e *ErrorDB) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func ParseError(dberr string) error {
	var keyErrors = map[string]string{
		"users_username_key": "This name is occupied",
		"users_email_key":    "This email is already occupied",
		"chats_pkey":         "This chat name is already occupied",
	}

	splitErr := strings.Split(dberr, " ")
	key := strings.Trim(splitErr[len(splitErr)-1], "\"")
	el, ok := keyErrors[key]

	if ok {
		return &ErrorDB{el}
	} else {
		return &ErrorDB{"Unknown database error"}
	}
}
