package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benny66/go-eventbus"
)

// 自定义事件
type FilterEvent struct {
	Category string
	Priority int
	Message  string
}

// 实现 Event 接口
func (e *FilterEvent) Topic() string {
	return "filter-topic"
}

// 过滤处理器 - 只处理特定类别的事件
type CategoryFilterHandler struct {
	category string
}

// 实现 Handler 接口
func (h *CategoryFilterHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*FilterEvent); ok {
		// 只处理指定类别的事件
		if e.Category == h.category {
			fmt.Printf("[类别过滤器] 处理 '%s' 类别事件: %s\n", e.Category, e.Message)
			return nil
		}
		// 忽略其他类别的事件
		fmt.Printf("[类别过滤器] 忽略 '%s' 类别事件\n", e.Category)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

// 优先级过滤处理器 - 只处理高优先级事件
type PriorityFilterHandler struct {
	minPriority int
}

// 实现 Handler 接口
func (h *PriorityFilterHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if e, ok := event.(*FilterEvent); ok {
		// 只处理高于最小优先级的事件
		if e.Priority >= h.minPriority {
			fmt.Printf("[优先级过滤器] 处理优先级 %d 的事件: %s\n", e.Priority, e.Message)
			return nil
		}
		// 忽略低优先级事件
		fmt.Printf("[优先级过滤器] 忽略优先级 %d 的事件\n", e.Priority)
		return nil
	}
	return fmt.Errorf("无效的事件类型")
}

func main() {
	// 创建事件总线
	eb := eventbus.NewEventBus()

	// 注册不同的过滤处理器
	eb.Subscribe("filter-topic", &CategoryFilterHandler{category: "system"})
	eb.Subscribe("filter-topic", &PriorityFilterHandler{minPriority: 5})

	// 发布不同类别和优先级的事件
	events := []FilterEvent{
		{Category: "system", Priority: 3, Message: "系统低优先级消息"},
		{Category: "user", Priority: 7, Message: "用户高优先级消息"},
		{Category: "system", Priority: 8, Message: "系统高优先级消息"},
		{Category: "user", Priority: 2, Message: "用户低优先级消息"},
	}

	for _, evt := range events {
		fmt.Printf("\n发布事件: 类别=%s, 优先级=%d\n", evt.Category, evt.Priority)
		eb.Publish(context.Background(), &evt)
		time.Sleep(time.Millisecond * 500) // 间隔发布，便于观察
	}

	// 等待事件处理完成
	time.Sleep(time.Second)

	// 关闭事件总线
	eb.Close()
}