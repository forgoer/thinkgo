package formatter

import "github.com/forgoer/thinkgo/log/record"

type Formatter interface {
	Format(r record.Record) string
	FormatBatch(rs []record.Record) string
}
