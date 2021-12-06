package model

import "encoding/json"

type Register struct {
	Username  string
	UserId    string
	NewClient bool
}

func (r Register) GetBytes() ([]byte, error) {
	return json.Marshal(r)
}
