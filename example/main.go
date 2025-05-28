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
// 日志处理器 - 专门记录事件日志
type LogHandler struct{}

// 实现 Handler 接口
func (h *LogHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*MyEvent); ok {
		fmt.Printf("[日志处理器] 记录事件: %s\n", e.Message)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

// 统计处理器 - 专门统计事件数据
type StatsHandler struct{}

// 实现 Handler 接口
func (h *StatsHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*MyEvent); ok {
		fmt.Printf("[统计处理器] 统计事件: 消息长度=%d\n", len(e.Message))
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}
func main() {
	// 创建事件总线
	// 自定义最大重试次数
	eb := eventbus.NewEventBus(eventbus.WithMaxRetries(5), eventbus.WithLogger(&MyLogger{}))

	// 多订阅者模式示例
	// 为同一个主题注册多个不同的处理器
	
	// 订阅者1: 普通处理器
	eb.Subscribe("my-topic", &MyHandler{})
	
	// 订阅者2: 日志处理器
	eb.Subscribe("my-topic", &LogHandler{})
	
	// 订阅者3: 统计处理器
	eb.Subscribe("my-topic", &StatsHandler{})

	// 发布一个事件，所有订阅者都会收到并处理
	eb.Publish(context.Background(), &MyEvent{Message: "Hello, EventBus!"})

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}
