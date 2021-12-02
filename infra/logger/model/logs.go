package model

import (
	"time"
)

type Log struct {
	Level   string    `json:"severity"`
	Message string    `json:"entity"`
	Error   string    `json:"error"`
	Time    time.Time `json:"timestamp"`
}
