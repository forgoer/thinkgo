package formatter

import (
	"encoding/json"
	"github.com/forgoer/thinkgo/log/record"
)

type JsonFormatter struct {
}

func NewJsonFormatter() *JsonFormatter {
	j := &JsonFormatter{}
	return j
}

func (f *JsonFormatter) Format(r record.Record) string {
	normalized := make(map[string]interface{})

	normalized["message"] = r.Message
	normalized["level"] = r.Level
	normalized["level_name"] = r.LevelName
	normalized["channel"] = r.Channel
	normalized["datetime"] = r.Datetime.Local().Format("2006-01-02 15:04:05.000")

	output, _ := json.Marshal(normalized)
	return string(output) + "\n"
}

func (f *JsonFormatter) FormatBatch(rs []record.Record) string {
	message := ""
	for _, r := range rs {
		message = message + f.Format(r)
	}
	return message
}
