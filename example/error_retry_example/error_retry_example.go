package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type RetryEvent struct {
	ID      int
	Message string
}

// 实现 Event 接口
func (e *RetryEvent) Topic() string {
	return "retry-topic"
}

// 会失败的处理器
type FailingHandler struct {
	failCount int
	maxFails  int
}

// 创建一个新的FailingHandler
func NewFailingHandler(maxFails int) *FailingHandler {
	return &FailingHandler{failCount: 0, maxFails: maxFails}
}

// 实现 Handler 接口
func (h *FailingHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*RetryEvent); ok {
		h.failCount++
		fmt.Printf("尝试处理事件 ID: %d, 尝试次数: %d\n", e.ID, h.failCount)

		// 模拟前几次处理失败，最后一次成功
		if h.failCount <= h.maxFails {
			return fmt.Errorf("处理失败，将重试 (尝试 %d/%d)", h.failCount, h.maxFails+1)
		}

		fmt.Printf("成功处理事件: %s (在 %d 次尝试后)\n", e.Message, h.failCount)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

func main() {
	// 创建事件总线，设置最大重试次数为3
	eb := eventbus.NewEventBus(
		eventbus.WithMaxRetries(3),
	)

	// 创建一个会失败2次然后成功的处理器
	failingHandler := NewFailingHandler(2)

	// 订阅事件
	eb.Subscribe("retry-topic", failingHandler)

	// 发布事件
	eb.Publish(context.Background(), &RetryEvent{
		ID:      1001,
		Message: "这是一个需要重试的事件",
	})

	// 等待事件处理完成（包括重试）
	time.Sleep(time.Second * 2)

	// 关闭事件总线
	eb.Close()
}
