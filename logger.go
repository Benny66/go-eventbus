package eventbus

import (
	"log"
)

// Logger 定义日志接口
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// defaultLogger 默认日志实现
type defaultLogger struct{}

// Infof 实现 Logger 接口的 Infof 方法
func (l *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// Errorf 实现 Logger 接口的 Errorf 方法
func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}
