package messages

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"udp-chat/infra/logger"
	"udp-chat/infra/redis"
	"udp-chat/internal/app/chat/messages/model"
)

func TestStore(t *testing.T) {
	msgObj := model.Message{
		Id:       123,
		UserId:   "123",
		Username: "Felipe",
		Text:     "Testing",
		Date:     time.Now(),
	}
	b, _ := json.Marshal(msgObj)
	log := new(logger.Mock)
	cache := new(redis.MockCache)
	cache.On("Get", CacheObjName).Return([]byte{}, errors.New("mock error")).Once()
	cache.On("Set", CacheObjName).Return(nil).Once()
	messageService := NewMessage(cache, log, 20)
	resp, err := messageService.Store(string(b))

	assert.NotNil(t, resp)
	assert.Nil(t, err)
}

func TestStoreWithUnmarshalError(t *testing.T) {
	msgObj := []model.Message{
		{
			Id:       123,
			UserId:   "123",
			Username: "Felipe",
			Text:     "Testing",
			Date:     time.Now(),
		},
	}
	b, _ := json.Marshal(msgObj)
	log := new(logger.Mock)
	cache := new(redis.MockCache)
	cache.On("Get", CacheObjName).Return(b, nil)

	messageService := NewMessage(cache, log, 20)
	resp, err := messageService.Store(string(b))

	assert.Nil(t, resp)
	assert.NotNil(t, err)
}
