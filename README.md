# Go EventBus

一个简单、高效的Go语言事件总线库，支持异步事件处理、自动重试和优雅关闭。

## 特性

- 基于主题的发布/订阅模式
- 异步事件处理
- 自动重试机制
- 优雅关闭
- 可自定义日志接口
- 线程安全

## 安装

```bash
go get github.com/benny66/go-eventbus

```

### example/main.go

```go
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

func main() {
	// 创建事件总线
	// 自定义channel容量
	eb := eventbus.NewEventBus(eventbus.WithChannelSize(100))

	// 订阅事件
	eb.Subscribe("my-topic", &MyHandler{})

	// 发布事件
	eb.Publish(context.Background(), &MyEvent{Message: "Hello, EventBus!"})

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}
```