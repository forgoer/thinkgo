package log

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/forgoer/thinkgo/log/handler"
	"github.com/forgoer/thinkgo/log/record"
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

func TestLogWithFileHandler(t *testing.T) {

	filename := path.Join(os.TempDir(), "thinkgo.log")

	h := handler.NewFileHandler(filename, record.INFO)

	l := NewLogger("testing", record.INFO)
	l.PushHandler(h)

	filename = h.GetTimedFilename()

	os.Remove(filename)

	message := "Log write to file"

	l.Debug(message)

	_, err := ioutil.ReadFile(filename)
	if err == nil {
		t.Error(errors.New("test FileHandler error"))
	}


	h.SetLevel(record.DEBUG)
	l = NewLogger("testing", record.DEBUG)
	l.PushHandler(h)
	l.Debug(message)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(errors.New("test FileHandler error"))
	}
	content := string(b)

	if !strings.Contains(content, message) {
		t.Error("test FileHandler error")
	}

}
