// Code generated by jsonenums -type=Kind; DO NOT EDIT.

package main

import (
	"encoding/json"
	"fmt"
)

var (
	_KindNameToValue = map[string]Kind{
		"NEW_MESSAGE":      NEW_MESSAGE,
		"NEW_CHAT":         NEW_CHAT,
		"SEARCH_QUERY":     SEARCH_QUERY,
		"NEW_PRIVATE_CHAT": NEW_PRIVATE_CHAT,
		"JOIN_CHAT":        JOIN_CHAT,
		"LOAD_MESSAGES":    LOAD_MESSAGES,
		"LOAD_SUBS":        LOAD_SUBS,
		"NEW_GROUP_CHAT":   NEW_GROUP_CHAT,
	}

	_KindValueToName = map[Kind]string{
		NEW_MESSAGE:      "NEW_MESSAGE",
		NEW_CHAT:         "NEW_CHAT",
		SEARCH_QUERY:     "SEARCH_QUERY",
		NEW_PRIVATE_CHAT: "NEW_PRIVATE_CHAT",
		JOIN_CHAT:        "JOIN_CHAT",
		LOAD_MESSAGES:    "LOAD_MESSAGES",
		LOAD_SUBS:        "LOAD_SUBS",
		NEW_GROUP_CHAT:   "NEW_GROUP_CHAT",
	}
)

func init() {
	var v Kind
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_KindNameToValue = map[string]Kind{
			interface{}(NEW_MESSAGE).(fmt.Stringer).String():      NEW_MESSAGE,
			interface{}(NEW_CHAT).(fmt.Stringer).String():         NEW_CHAT,
			interface{}(SEARCH_QUERY).(fmt.Stringer).String():     SEARCH_QUERY,
			interface{}(NEW_PRIVATE_CHAT).(fmt.Stringer).String(): NEW_PRIVATE_CHAT,
			interface{}(JOIN_CHAT).(fmt.Stringer).String():        JOIN_CHAT,
			interface{}(LOAD_MESSAGES).(fmt.Stringer).String():    LOAD_MESSAGES,
			interface{}(LOAD_SUBS).(fmt.Stringer).String():        LOAD_SUBS,
			interface{}(NEW_GROUP_CHAT).(fmt.Stringer).String():   NEW_GROUP_CHAT,
		}
	}
}

// MarshalJSON is generated so Kind satisfies json.Marshaler.
func (r Kind) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _KindValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid Kind: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so Kind satisfies json.Unmarshaler.
func (r *Kind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Kind should be a string, got %s", data)
	}
	v, ok := _KindNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid Kind %q", s)
	}
	*r = v
	return nil
}
