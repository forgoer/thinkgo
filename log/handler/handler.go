package handler

import (
	"github.com/forgoer/thinkgo/log/formatter"
	"github.com/forgoer/thinkgo/log/record"
)

type Handler struct {
	formatter formatter.Formatter
	level    record.Level
}

// GetFormatter Gets the formatter.
func (h *Handler) GetFormatter() formatter.Formatter {
	if h.formatter == nil {
		h.formatter = h.getDefaultFormatter()
	}
	return h.formatter
}

// SetFormatter Sets the formatter.
func (h *Handler) SetFormatter(f formatter.Formatter) *Handler {
	h.formatter = f
	return h
}

func (h *Handler) getDefaultFormatter() formatter.Formatter {
	return formatter.NewLineFormatter()
}
