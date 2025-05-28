package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type DynamicEvent struct {
	ID      int
	Message string
}

// 实现 Event 接口
func (e *DynamicEvent) Topic() string {
	return "dynamic-topic"
}

// 动态处理器
type DynamicHandler struct {
	name string
}

// 实现 Handler 接口
func (h *DynamicHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*DynamicEvent); ok {
		fmt.Printf("[%s] 处理事件 ID: %d, 消息: %s\n", h.name, e.ID, e.Message)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

func main() {
	// 创建事件总线
	eb := eventbus.NewEventBus()

	// 创建多个处理器
	handler1 := &DynamicHandler{name: "处理器1"}
	handler2 := &DynamicHandler{name: "处理器2"}
	handler3 := &DynamicHandler{name: "处理器3"}

	// 初始只订阅处理器1
	fmt.Println("阶段1: 只有处理器1订阅事件")
	eb.Subscribe("dynamic-topic", handler1)

	// 发布事件1
	eb.Publish(context.Background(), &DynamicEvent{ID: 1, Message: "第一条消息"})
	time.Sleep(time.Millisecond * 500)

	// 添加处理器2
	fmt.Println("\n阶段2: 添加处理器2")
	eb.Subscribe("dynamic-topic", handler2)

	// 发布事件2
	eb.Publish(context.Background(), &DynamicEvent{ID: 2, Message: "第二条消息"})
	time.Sleep(time.Millisecond * 500)

	// 添加处理器3并移除处理器1
	fmt.Println("\n阶段3: 添加处理器3，移除处理器1")
	eb.Subscribe("dynamic-topic", handler3)
	eb.Unsubscribe("dynamic-topic", handler1)

	// 发布事件3
	eb.Publish(context.Background(), &DynamicEvent{ID: 3, Message: "第三条消息"})
	time.Sleep(time.Millisecond * 500)

	// 移除所有处理器
	fmt.Println("\n阶段4: 移除所有处理器")
	eb.Unsubscribe("dynamic-topic", handler2)
	eb.Unsubscribe("dynamic-topic", handler3)

	// 发布事件4 (不会被处理，因为没有订阅者)
	eb.Publish(context.Background(), &DynamicEvent{ID: 4, Message: "第四条消息"})
	fmt.Println("发布了第四条消息，但没有处理器接收")

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}