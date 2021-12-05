package entity

import "time"

type Message struct {
	Id       string
	UserId   string
	Username string
	Text     string
	Date     time.Time
}
