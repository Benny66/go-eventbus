package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type ConfigEvent struct {
	Message string
}

// 实现 Event 接口
func (e *ConfigEvent) Topic() string {
	return "config-topic"
}

// 自定义处理器
type ConfigHandler struct{}

// 实现 Handler 接口
func (h *ConfigHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*ConfigEvent); ok {
		fmt.Printf("处理配置事件: %s\n", e.Message)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

// 自定义日志记录器
type CustomLogger struct{}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[自定义INFO] "+format+"\n", args...)
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[自定义ERROR] "+format+"\n", args...)
}

func main() {
	// 创建事件总线，自定义配置
	eb := eventbus.NewEventBus(
		// 设置最大重试次数为5
		eventbus.WithMaxRetries(5),
		// 设置事件通道容量为500
		eventbus.WithChannelSize(500),
		// 设置自定义日志记录器
		eventbus.WithLogger(&CustomLogger{}),
	)

	// 订阅事件
	eb.Subscribe("config-topic", &ConfigHandler{})

	// 发布事件
	eb.Publish(context.Background(), &ConfigEvent{Message: "测试自定义配置!"})

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}