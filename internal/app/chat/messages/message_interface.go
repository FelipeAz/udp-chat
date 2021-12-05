package messages

import "udp-chat/internal/app/chat/messages/model"

type MessageInterface interface {
	Store(msg string) (*model.Message, error)
	Get() ([]model.Message, error)
	Delete(id string) error
}
