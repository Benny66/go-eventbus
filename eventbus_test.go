package eventbus

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(ctx context.Context, event Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// 自定义日志记录器
type MockLogger struct{}

func (l *MockLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[自定义INFO] "+format+"\n", args...)
}

func (l *MockLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[自定义ERROR] "+format+"\n", args...)
}

func TestNewEventBus(t *testing.T) {
	t.Run("Default Configuration", func(t *testing.T) {
		eb := NewEventBus()
		assert.NotNil(t, eb.handlers)
		assert.NotNil(t, eb.eventChans)
		assert.Equal(t, 3, eb.maxRetries)
		assert.Equal(t, 1000, eb.channelSize)
		assert.NotNil(t, eb.done)
		assert.NotNil(t, eb.logger)
	})

	t.Run("WithMaxRetries", func(t *testing.T) {
		eb := NewEventBus(WithMaxRetries(5))
		assert.Equal(t, 5, eb.maxRetries)
	})

	t.Run("WithChannelSize", func(t *testing.T) {
		eb := NewEventBus(WithChannelSize(500))
		assert.Equal(t, 500, eb.channelSize)
	})

	t.Run("WithLogger", func(t *testing.T) {
		logger := new(MockLogger)
		eb := NewEventBus(WithLogger(logger))
		assert.Same(t, logger, eb.logger)
	})

	t.Run("WithChannelSize_Zero", func(t *testing.T) {
		eb := NewEventBus(WithChannelSize(0))
		assert.Equal(t, 1000, eb.channelSize)
	})

	t.Run("WithChannelSize_Negative", func(t *testing.T) {
		eb := NewEventBus(WithChannelSize(-100))
		assert.Equal(t, 1000, eb.channelSize)
	})
}
