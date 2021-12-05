package model

import (
	"encoding/json"
	"time"
)

type Message struct {
	Id       int
	UserId   string
	Username string
	Text     string
	Date     time.Time
}

func NewMessage(id int, username, userId, text string) Message {
	return Message{
		Id:       id,
		UserId:   userId,
		Username: username,
		Text:     text,
		Date:     time.Now().Local(),
	}
}

func (m Message) GetDateFormated(format string) string {
	return m.Date.Format(format)
}

func (m Message) ToBytes() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}
