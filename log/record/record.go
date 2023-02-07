package record

import (
	"strings"
	"time"
)

type Level int

const (
	// DEBUG Detailed debug information
	DEBUG Level = 100
	// INFO Interesting events
	INFO Level = 200
	// NOTICE Uncommon events
	NOTICE Level = 250
	// WARNING Exceptional occurrences that are not errors
	WARNING Level = 300
	// ERROR Runtime errors
	ERROR Level = 400
	// CRITICAL Critical conditions
	CRITICAL Level = 500
	// ALERT Action must be taken immediately
	ALERT Level = 550
	// EMERGENCY Urgent alert.
	EMERGENCY Level = 600
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

// GetLevel Parse the string level into a Level constant.
func GetLevel(levelKey string) Level {
	for level, s := range levels {
		if strings.ToUpper(s) == s {
			return level
		}
	}
	return DEBUG
}
