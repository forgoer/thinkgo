package filesystem

import (
	"os"
	"path"
	"io/ioutil"
	"sync"
	"time"
	"path/filepath"
)

var lock sync.RWMutex

func Exists(p ...string) (bool, error) {
	_, err := os.Stat(path.Join(p...))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Get(path string) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()

	var b []byte

	f, err := os.OpenFile(path, os.O_RDWR, 0600)
	defer f.Close()

	if err != nil {
		return b, err
	}

	b, err = ioutil.ReadAll(f)
	if err != nil {
		return b, err
	}
	return b, nil
}

func Put(path string, data string) error {
	lock.Lock()
	defer lock.Unlock()

	dir := filepath.Dir(path)
	if ok, _ := Exists(dir); !ok {
		err := os.MkdirAll(dir, 0600)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	defer f.Close()

	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	return err
}

func ModTime(path string) (time.Time, error) {
	var modTime time.Time
	fileInfo, err := os.Stat(path)
	if err != nil {
		return modTime, err
	}
	return fileInfo.ModTime(), nil
}
