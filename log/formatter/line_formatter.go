package formatter

import "github.com/forgoer/thinkgo/log/record"

type LineFormatter struct {
}

func NewLineFormatter() *LineFormatter {
	f := &LineFormatter{}
	return f
}

func (f *LineFormatter) Format(r record.Record) string {
	return "[" + r.Datetime.Local().Format("2006-01-02 15:04:05.000") + "]" + r.Channel + "." + r.LevelName + ": " + r.Message + "\n"
}

func (f *LineFormatter) FormatBatch(rs []record.Record) string {
	message := ""
	for _, r := range rs {
		message = message + f.Format(r)
	}
	return message
}
