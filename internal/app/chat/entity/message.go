package entity

import "time"

type Message struct {
	Id   string
	Text string
	Date time.Time
}
