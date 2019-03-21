package log

import (
	"container/list"
	"errors"
	"fmt"
	"time"

	"github.com/thinkoner/thinkgo/log/handler"
	"github.com/thinkoner/thinkgo/log/record"
)

type Handler interface {
	IsHandling(r record.Record) bool
	Handle(r record.Record) bool
}

type Logger struct {
	name     string
	level    record.Level
	handlers *list.List
}

// NewLogger New a Logger instance
func NewLogger(name string, level record.Level) *Logger {
	logger := &Logger{
		name:     name,
		handlers: list.New(),
		level:    level,
	}
	return logger
}

// GetName Gets the name
func (logger *Logger) GetName(name string) string {
	return logger.name
}

// SetName Sets the name
func (logger *Logger) SetName(name string) *Logger {
	logger.name = name
	return logger
}

// PushHandler Pushes a handler on to the stack.
func (logger *Logger) PushHandler(handler Handler) *Logger {
	logger.handlers.PushFront(handler)
	return logger
}

// PopHandler Pops a handler from the stack
func (logger *Logger) PopHandler() {
	front := logger.handlers.Front()
	if front != nil {
		logger.handlers.Remove(front)
	}
}

// SetHandlers Set handlers, replacing all existing ones.
func (logger *Logger) SetHandlers(handlers []Handler) *Logger {
	count := len(handlers)
	for i := count - 1; i >= 0; i = i - 1 {
		logger.PushHandler(handlers[i])
	}
	return logger
}

// GetHandlers Returns a Handler slice
func (logger *Logger) GetHandlers() []Handler {
	var handler []Handler
	for e := logger.handlers.Front(); e != nil; e = e.Next() {
		handler = append(handler, e.Value.(Handler))
	}
	return handler
}

// AddRecord Adds a log record.
func (logger *Logger) AddRecord(level record.Level, format string, v ...interface{}) (bool, error) {
	if logger.handlers.Len() == 0 {
		logger.PushHandler(handler.NewConsoleHandler(logger.level))
	}

	levelName, err := GetLevelName(level)
	if err != nil {
		return false, err
	}

	handlerKey := false
	for e := logger.handlers.Front(); e != nil; e = e.Next() {
		h := e.Value.(Handler)
		if h.IsHandling(record.Record{Level: level}) {
			handlerKey = true
			break
		}
	}
	if !handlerKey {
		return false, nil
	}

	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}

	r := record.Record{
		Level:     level,
		Message:   format,
		LevelName: levelName,
		Channel:   logger.name,
		Datetime:  time.Now(),
	}

	for e := logger.handlers.Front(); e != nil; e = e.Next() {
		h := e.Value.(Handler)
		if h.Handle(r) {
			break
		}
	}

	return true, nil
}

// Debug Adds a log record at the DEBUG level.
func (logger *Logger) Debug(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.DEBUG, format, v...)
}

// Info Adds a log record at the INFO level.
func (logger *Logger) Info(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.INFO, format, v...)
}

// Notice Adds a log record at the NOTICE level.
func (logger *Logger) Notice(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.NOTICE, format, v...)
}

// Warn Adds a log record at the WARNING level.
func (logger *Logger) Warn(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.WARNING, format, v...)
}

// Error Adds a log record at the ERROR level.
func (logger *Logger) Error(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.ERROR, format, v...)
}

// Crit Adds a log record at the CRITICAL level.
func (logger *Logger) Crit(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.CRITICAL, format, v...)
}

// Alert Adds a log record at the ALERT level.
func (logger *Logger) Alert(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.ALERT, format, v...)
}

// Emerg Adds a log record at the EMERGENCY level.
func (logger *Logger) Emerg(format string, v ...interface{}) (bool, error) {
	return logger.AddRecord(record.EMERGENCY, format, v...)
}

// Gets the name of the logging level.
func GetLevelName(level record.Level) (string, error) {
	levels := record.GetLevels()
	l, ok := levels[level]
	if !ok {
		return l, errors.New(fmt.Sprintf("Level %d is not defined", level))
	}
	return l, nil
}
