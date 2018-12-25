package route

import (
	"net/http"

	"github.com/thinkoner/thinkgo/filesystem"
)

type staticHandle struct {
	fileServer http.Handler
	fs         http.FileSystem
}

func NewStaticHandle(root string) http.Handler {
	fs := filesystem.NewFileFileSystem(root, false)
	return &staticHandle{
		fileServer: http.FileServer(fs),
		fs:         fs,
	}
}

func (s *staticHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.fileServer.ServeHTTP(w, r)
}
