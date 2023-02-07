package handler

import (
	"github.com/forgoer/thinkgo/log/record"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestRotateHandler(t *testing.T) {
	filename := path.Join(os.TempDir(), "thinkgo.log")

	h := NewRotateHandler(filename, record.DEBUG)
	filename = h.GetFilename()

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
	t.Log(filename, content)
}
