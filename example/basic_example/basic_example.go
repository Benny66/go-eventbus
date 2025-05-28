package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type BasicEvent struct {
	Message string
}

// 实现 Event 接口
func (e *BasicEvent) Topic() string {
	return "basic-topic"
}

// 自定义处理器
type BasicHandler struct{}

// 实现 Handler 接口
func (h *BasicHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*BasicEvent); ok {
		fmt.Printf("处理基础事件: %s\n", e.Message)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

func main() {
	// 创建事件总线
	eb := eventbus.NewEventBus()

	// 订阅事件
	eb.Subscribe("basic-topic", &BasicHandler{})

	// 发布事件
	eb.Publish(context.Background(), &BasicEvent{Message: "这是一个基础示例!"})

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}