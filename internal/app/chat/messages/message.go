package messages

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
	"udp-chat/internal/app/chat/messages/model"
	error_messages "udp-chat/internal/app/chat/server/constants"
	database "udp-chat/internal/cache"
	"udp-chat/internal/logger"
)

const (
	CacheObjName = "CHAT"
)

type Message struct {
	Cache  database.CacheInterface
	Logger logger.LogInterface
	Size   int
}

func NewMessage(cache database.CacheInterface, log logger.LogInterface, maxSize int) Message {
	return Message{
		Cache:  cache,
		Logger: log,
		Size:   maxSize,
	}
}

func (m Message) Store(msgObj *model.Message) error {
	messages, err := m.Get()
	if err != nil {
		err = errors.Wrap(err, error_messages.FailedToGetMessagesFromChat)
		m.Logger.Error(err)
	}

	msgObj.Date = time.Now()
	messages = m.addMessageToQueue(messages, *msgObj)

	b, err := json.Marshal(messages)
	if err != nil {
		m.Logger.Error(err)
		return err
	}

	err = m.Cache.Set(CacheObjName, b)
	if err != nil {
		m.Logger.Error(err)
		return err
	}

	return nil
}

func (m Message) Get() ([]model.Message, error) {
	var messages []model.Message
	b, err := m.Cache.Get(CacheObjName)
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

func (m Message) addMessageToQueue(queue []model.Message, msg model.Message) []model.Message {
	var newQueue []model.Message
	if len(queue) < m.Size {
		newQueue = append(queue, msg)
		return newQueue
	}

	newQueue = append(queue[1:], msg)
	return newQueue
}

func (m Message) Delete(id, userId string) error {
	return nil
}

func (m Message) UnmarshalMessage(b []byte) (model.Message, error) {
	var msg model.Message
	err := json.Unmarshal(b, &msg)
	return msg, err
}
