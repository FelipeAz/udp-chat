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
	Id        int
	UserId    string
	Username  string
	Text      string
	NewClient bool
	Date      time.Time
}

func NewMessage(id int, username, userId, text string) Message {
	return Message{
		Id:        id,
		UserId:    userId,
		Username:  username,
		Text:      text,
		NewClient: false,
		Date:      time.Now().UTC(),
	}
}

func (m Message) GetMessageFormated() string {
	if m.NewClient {
		return fmt.Sprintf("%s Joined the chat", m.Username)
	}
	date := m.GetDateFormated()
	return fmt.Sprintf("%s %s: %s", date, m.Username, m.Text)
}

func (m Message) GetDateFormated() string {
	return m.Date.Local().Format(dateFormat)
}

func (m Message) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}
