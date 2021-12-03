package messages

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
	"udp-chat/internal/app/chat/entity"
	error_messages "udp-chat/internal/app/chat/server/constants"
	database "udp-chat/internal/cache"
	"udp-chat/internal/logger"
)

type Message struct {
	Cache  database.CacheInterface
	Logger logger.LogInterface
}

func NewMessage(cache database.CacheInterface, log logger.LogInterface) Message {
	return Message{
		Cache:  cache,
		Logger: log,
	}
}

func (m Message) StoreMessage(msg string) (string, error) {
	messages, err := m.GetMessages()
	if err != nil {
		err = errors.Wrap(err, error_messages.FailedToGetMessagesFromChat)
		m.Logger.Error(err)
	}
	msgObj := entity.Message{
		Id:   uuid.NewString(),
		Text: msg,
		Date: time.Now(),
	}

	messages = append(messages, msgObj)

	b, err := json.Marshal(messages)
	if err != nil {
		m.Logger.Error(err)
		return "", err
	}

	err = m.Cache.Set("CHAT", b)
	if err != nil {
		m.Logger.Error(err)
		return "", err
	}

	return msgObj.Id, nil
}

func (m Message) GetMessages() ([]entity.Message, error) {
	var messages []entity.Message
	b, err := m.Cache.Get("CHAT")
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal(b, &messages)
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}

	return messages, nil
}

func (m Message) DeleteMessage(id string) error {
	return nil
}
