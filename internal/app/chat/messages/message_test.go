package messages

import (
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
	log := new(logger.Mock)
	cache := new(redis.MockCache)
	cache.On("Get", CacheObjName).Return([]byte{}, errors.New("mock error")).Once()
	cache.On("Set", CacheObjName).Return(nil).Once()
	messageService := NewMessage(cache, log, 20)
	err := messageService.Store(&msgObj)

	assert.Nil(t, err)
}
