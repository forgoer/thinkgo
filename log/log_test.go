package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Debug("log with Debug")
	Info("log with Info")
	Notice("log with Notice")
	Warn("log with Warn")
	Error("log with Error")
	Crit("log with Crit")
	Alert("log with Alert")
	Emerg("log with Emerg")
}
