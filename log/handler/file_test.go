package handler

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/forgoer/thinkgo/log/record"
)

func TestNewFileHandler(t *testing.T) {
	filename := path.Join(os.TempDir(), "thinkgo.log")

	h := NewFileHandler(filename, record.DEBUG)
	filename = h.GetTimedFilename()

	os.Remove(filename)

	message := "Log write to file"
	r := getRecord(message)
	h.Handle(r)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	content := string(b)

	if !strings.Contains(content, message) {
		t.Error("test FileHandler error")
	}

}

func getRecord(message string) record.Record {
	return record.Record{
		Level:     200,
		Message:   message,
		LevelName: "INFO",
		Channel:   "testing",
		Datetime:  time.Now(),
	}
}
