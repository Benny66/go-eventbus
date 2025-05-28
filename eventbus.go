package eventbus

import (
	"context"
	"sync"
)

// Event 事件接口
type Event interface {
	Topic() string
}

// Handler 事件处理器接口
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// eventTask 事件任务
type eventTask struct {
	ctx     context.Context
	event   Event
	handler Handler
}

// EventBus 事件总线
type EventBus struct {
	handlers    map[string][]Handler
	mutex       sync.RWMutex
	eventChans  map[string]chan eventTask
	maxRetries  int
	channelSize int // 新增字段：channel容量
	// 用于关闭和同步
	done chan struct{}  // 用于关闭事件总线
	wg   sync.WaitGroup // 等待所有 goroutine 结束
	// 日志接口
	logger Logger
}

// EventBusOption 配置选项函数类型
type EventBusOption func(*EventBus)

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) EventBusOption {
	return func(eb *EventBus) {
		eb.maxRetries = maxRetries
	}
}

// WithChannelSize 设置事件队列channel容量
func WithChannelSize(size int) EventBusOption {
	return func(eb *EventBus) {
		if size > 0 {
			eb.channelSize = size
		}
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger Logger) EventBusOption {
	return func(eb *EventBus) {
		eb.logger = logger
	}
}

// NewEventBus 创建一个新的事件总线
// 参数 options: 可变参数，用于自定义事件总线的配置选项
// 返回值: 初始化完成的EventBus实例
func NewEventBus(options ...EventBusOption) *EventBus {
	// 初始化EventBus实例，设置默认值
	eb := &EventBus{
		// handlers: 存储事件处理器的映射，key为事件主题，value为该主题的处理器列表
		// 允许一个主题有多个处理器，实现多订阅者模式
		handlers: make(map[string][]Handler),

		// eventChans: 存储事件通道的映射，key为事件主题，value为该主题的事件队列
		// 用于异步处理事件，实现事件的缓冲和并发处理
		eventChans: make(map[string]chan eventTask),

		// maxRetries: 事件处理失败时的最大重试次数
		// 默认为3次，可通过WithMaxRetries选项修改
		maxRetries: 3,

		// channelSize: 事件队列的缓冲区大小
		// 默认为1000，可通过WithChannelSize选项修改
		// 较大的缓冲区可以处理突发的大量事件，但会占用更多内存
		channelSize: 1000,

		// done: 用于通知所有goroutine退出的信号通道
		// 在Close方法中关闭，用于实现优雅关闭
		done: make(chan struct{}),

		// logger: 日志记录器，默认使用内置的defaultLogger
		// 可通过WithLogger选项自定义
		logger: &defaultLogger{},
	}

	// 应用自定义选项
	// 遍历所有传入的选项函数，并依次应用到EventBus实例
	for _, option := range options {
		option(eb)
	}

	return eb
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(topic string, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if _, ok := eb.handlers[topic]; !ok {
		eb.handlers[topic] = make([]Handler, 0)
		eb.eventChans[topic] = make(chan eventTask, eb.channelSize) // 使用配置的channel容量
		eb.wg.Add(1)
		go eb.processEvents(topic)
	}
	eb.handlers[topic] = append(eb.handlers[topic], handler)
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(topic string, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	handlers, exists := eb.handlers[topic]
	if !exists {
		return
	}

	// 移除指定的 handler
	for i, h := range handlers {
		if h == handler {
			eb.handlers[topic] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	// 如果没有订阅者了，清理资源
	if len(eb.handlers[topic]) == 0 {
		close(eb.eventChans[topic])
		delete(eb.handlers, topic)
		delete(eb.eventChans, topic)
	}
}

// processEvents 处理事件队列
func (eb *EventBus) processEvents(topic string) {
	defer eb.wg.Done()

	eventChan := eb.eventChans[topic]
	for {
		// 打印当前队列长度
		eb.logger.Infof("当前事件队列长度 - topic: %s, length: %d", topic, len(eventChan))

		select {
		case task, ok := <-eventChan:
			if !ok {
				return
			}
			var err error
			for i := 0; i < eb.maxRetries; i++ {
				ctx := context.Background() // 新建一个 context, 超时时间下层可以自己设置
				// 增加异常处理和容错机制
				err = eb.safeHandle(ctx, task.event, task.handler)
				if err == nil {
					break
				}
				eb.logger.Errorf("event topic %s handle failed: %s, retrying %d times...",
					task.event.Topic(), err.Error(), i+1)
				if i < eb.maxRetries-1 {
					continue
				}
			}
		case <-eb.done:
			return
		}
	}
}

// safeHandle 安全处理事件
func (eb *EventBus) safeHandle(ctx context.Context, event Event, handler Handler) error {
	defer func() {
		if r := recover(); r != nil {
			eb.logger.Errorf("event topic %s handle panic: %v", event.Topic(), r)
		}
	}()

	return handler.Handle(ctx, event)
}

// Close 关闭事件总线
func (eb *EventBus) Close() {
	eb.mutex.Lock()
	// 关闭所有 channel
	for topic, ch := range eb.eventChans {
		close(ch)
		delete(eb.eventChans, topic)
		delete(eb.handlers, topic)
	}
	close(eb.done)
	eb.mutex.Unlock()

	// 等待所有 goroutine 结束
	eb.wg.Wait()
}

// Publish 发布事件
func (eb *EventBus) Publish(ctx context.Context, event Event) {
	eb.mutex.RLock()
	handlers, exists := eb.handlers[event.Topic()]
	eventChan, _ := eb.eventChans[event.Topic()]
	eb.mutex.RUnlock()

	if !exists {
		return
	}
	for _, handler := range handlers {
		eventChan <- eventTask{
			ctx:     ctx,
			event:   event,
			handler: handler,
		}
	}
}
