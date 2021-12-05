package model

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	dateFormat = "01/02/2006 03:04PM"
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
		Date:     time.Now().UTC(),
	}
}

func (m Message) GetMessageFormated() string {
	date := m.GetDateFormated()
	return fmt.Sprintf("%s %s: %s", date, m.Username, m.Text)
}

func (m Message) GetDateFormated() string {
	return m.Date.Local().Format(dateFormat)
}

func (m Message) ToBytes() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}
