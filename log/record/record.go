package record

import "time"

type Level int

const (
	// Detailed debug information
	DEBUG = 100
	// Interesting events
	INFO = 200
	// Uncommon events
	NOTICE = 250
	// Exceptional occurrences that are not errors
	WARNING = 300
	// Runtime errors
	ERROR = 400
	// Critical conditions
	CRITICAL = 500
	// Action must be taken immediately
	ALERT = 550
	// Urgent alert.
	EMERGENCY = 600
)

// Logging levels from syslog protocol defined in RFC 5424
var levels = map[Level]string{
	DEBUG:     "DEBUG",
	INFO:      "INFO",
	NOTICE:    "NOTICE",
	WARNING:   "WARNING",
	ERROR:     "ERROR",
	CRITICAL:  "CRITICAL",
	ALERT:     "ALERT",
	EMERGENCY: "EMERGENCY",
}

type Record struct {
	Level     Level
	Message   string
	LevelName string
	Channel   string
	Datetime  time.Time
	Formatted string
}

// GetLevels returns levels map
func GetLevels() map[Level]string {
	return levels
}
