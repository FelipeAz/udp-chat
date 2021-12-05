package messages

import "udp-chat/internal/app/chat/entity"

type MessageInterface interface {
	Store(msg string) (*entity.Message, error)
	Get() ([]entity.Message, error)
	Delete(id string) error
}
