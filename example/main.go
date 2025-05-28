package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type MyEvent struct {
	Message string
}

// 实现 Event 接口
func (e *MyEvent) Topic() string {
	return "my-topic"
}

// 自定义处理器
type MyHandler struct{}

// 实现 Handler 接口
func (h *MyHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*MyEvent); ok {
		fmt.Printf("处理事件: %s\n", e.Message)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

// 自定义日志记录器
type MyLogger struct{}

func (l *MyLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}
func main() {
	// 创建事件总线
	// 自定义最大重试次数
	eb := eventbus.NewEventBus(eventbus.WithMaxRetries(5), eventbus.WithLogger(&MyLogger{}))

	// 订阅事件
	eb.Subscribe("my-topic", &MyHandler{})

	// 发布事件
	eb.Publish(context.Background(), &MyEvent{Message: "Hello, EventBus!"})

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}
