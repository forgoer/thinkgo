package handler

import (
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/forgoer/thinkgo/log/record"
)

type FileHandler struct {
	Handler
	sync.Mutex
	level    record.Level
	bubble   bool
	filename string

	filenameFormat string
	dateFormat     string
	timedFilename  string
}

func NewFileHandler(filename string, level record.Level) *FileHandler {
	h := &FileHandler{
		level:          level,
		bubble:         true,
		filename:       filename,
		filenameFormat: "{filename}-{date}",
		dateFormat:     "2006-01-02",
	}
	// h.timedFilename = h.GetTimedFilename()
	return h
}

// IsHandling Checks whether the given record will be handled by this handler.
func (h *FileHandler) IsHandling(r record.Record) bool {
	return r.Level >= h.level
}

// Handle Handles a record.
func (h *FileHandler) Handle(r record.Record) bool {
	if !h.IsHandling(r) {
		return false
	}

	r.Formatted = h.GetFormatter().Format(r)

	h.write(r)

	return false == h.bubble
}

// SetLevel Sets minimum logging level at which this handler will be triggered.
func (h *FileHandler) SetLevel(level record.Level) {
	h.level = level
}

func (h *FileHandler) write(r record.Record) {
	h.Lock()
	defer h.Unlock()
	file, _ := os.OpenFile(h.GetTimedFilename(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer file.Close()
	file.Write([]byte(r.Formatted))
}

// GetTimedFilename Gets the Timed filename.
func (h *FileHandler) GetTimedFilename() string {
	dirname := path.Dir(h.filename)
	filename := path.Base(h.filename)
	fileExt := path.Ext(h.filename)
	filename = strings.TrimSuffix(filename, fileExt)

	timedFilename := strings.Replace(path.Join(dirname, h.filenameFormat), "{filename}", filename, -1)
	timedFilename = strings.Replace(timedFilename, "{date}", time.Now().Local().Format(h.dateFormat), -1)

	if len(fileExt) > 0 {
		timedFilename = timedFilename + fileExt
	}

	return timedFilename
}
