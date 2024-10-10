package errordb

import "fmt"

type ErrorDB struct {
	Message string
}

func (e *ErrorDB) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func ParseError(dberr string) error {
	chatUniqErr := `pq: duplicate key value violates unique constraint "chats_pkey"`
	switch dberr {
	case chatUniqErr:
		return &ErrorDB{"This chat name is occupied."}
	}
	return nil
}
