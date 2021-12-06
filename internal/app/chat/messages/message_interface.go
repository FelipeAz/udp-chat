package messages

import "udp-chat/internal/app/chat/messages/model"

type MessageInterface interface {
	Store(msgObj *model.Message) error
	Get() ([]model.Message, error)
	Delete(id, userId string) error
	UnmarshalMessage([]byte) (model.Message, error)
}
