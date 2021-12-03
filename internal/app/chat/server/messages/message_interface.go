package messages

import "udp-chat/internal/app/chat/entity"

type MessageInterface interface {
	StoreMessage(msg string) (string, error)
	GetMessages() ([]entity.Message, error)
	DeleteMessage(id string) error
}
