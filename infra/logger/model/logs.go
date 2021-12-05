package model

import (
	"time"
)

type Log struct {
	Level   string    `json:"severity"`
	Service string    `json:"service"`
	Message string    `json:"model"`
	Error   string    `json:"error"`
	Time    time.Time `json:"timestamp"`
}
