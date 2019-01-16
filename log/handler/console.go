package handler

import (
	"os"
	"sync"

	"github.com/thinkoner/thinkgo/log/formatter"
	"github.com/thinkoner/thinkgo/log/record"
)

type brush func(string) string

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = map[int]brush{
	record.EMERGENCY: newBrush("1;41"), // Emergency          Red background
	record.ALERT:     newBrush("1;35"), // Alert              purple
	record.CRITICAL:  newBrush("1;34"), // Critical           blue
	record.ERROR:     newBrush("1;31"), // Error              red
	record.WARNING:   newBrush("1;33"), // Warn               yellow
	record.INFO:      newBrush("1;36"), // Informational      sky blue
	record.DEBUG:     newBrush("1;32"), // Debug              green
	record.NOTICE:    newBrush("1;32"), // Trace              green
}

type ConsoleHandler struct {
	sync.Mutex
	level     int
	formatter formatter.Formatter
	bubble    bool
}

func NewConsoleHandler() *ConsoleHandler {
	return &ConsoleHandler{
		level:  record.DEBUG,
		bubble: true,
	}
}

// IsHandling Checks whether the given record will be handled by this handler.
func (h *ConsoleHandler) IsHandling(r record.Record) bool {
	return r.Level >= h.level
}

// Handle Handles a record.
func (h *ConsoleHandler) Handle(r record.Record) bool {
	if !h.IsHandling(r) {
		return false
	}

	r.Formatted = h.GetFormatter().Format(r)

	h.write(r)

	return false == h.bubble
}

func (h *ConsoleHandler) write(r record.Record) {
	h.Lock()
	defer h.Unlock()
	message := colors[r.Level](r.Formatted)
	os.Stdout.Write(append([]byte(message)))
}

// GetFormatter Gets the formatter.
func (h *ConsoleHandler) GetFormatter() formatter.Formatter {
	if h.formatter == nil {
		h.formatter = h.getDefaultFormatter()
	}
	return h.formatter
}

// SetFormatter Sets the formatter.
func (h *ConsoleHandler) SetFormatter(f formatter.Formatter) *ConsoleHandler {
	h.formatter = f
	return h
}

func (h *ConsoleHandler) getDefaultFormatter() formatter.Formatter {
	return formatter.NewLineFormatter()
}
