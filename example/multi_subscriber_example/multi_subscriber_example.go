package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type UserEvent struct {
	UserID   int
	Username string
	Action   string
}

// 实现 Event 接口
func (e *UserEvent) Topic() string {
	return "user-events"
}

// 日志处理器 - 记录用户事件
type LogHandler struct{}

// 实现 Handler 接口
func (h *LogHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*UserEvent); ok {
		fmt.Printf("[日志] 用户 %s (ID: %d) 执行了 %s 操作\n", 
			e.Username, e.UserID, e.Action)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

// 通知处理器 - 发送用户操作通知
type NotificationHandler struct{}

// 实现 Handler 接口
func (h *NotificationHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*UserEvent); ok {
		fmt.Printf("[通知] 发送通知: 用户 %s 刚刚执行了 %s 操作\n", 
			e.Username, e.Action)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

// 统计处理器 - 统计用户操作
type AnalyticsHandler struct{}

// 实现 Handler 接口
func (h *AnalyticsHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*UserEvent); ok {
		fmt.Printf("[统计] 记录用户行为: UserID=%d, Action=%s\n", 
			e.UserID, e.Action)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

func main() {
	// 创建事件总线
	eb := eventbus.NewEventBus()

	// 为同一个主题注册多个不同的处理器
	eb.Subscribe("user-events", &LogHandler{})
	eb.Subscribe("user-events", &NotificationHandler{})
	eb.Subscribe("user-events", &AnalyticsHandler{})

	// 发布一个事件，所有订阅者都会收到并处理
	eb.Publish(context.Background(), &UserEvent{
		UserID:   1001,
		Username: "张三",
		Action:   "登录",
	})

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}