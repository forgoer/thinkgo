package ctx

import (
	"io"
	"mime/multipart"
	"os"
	"path"
)

type File struct {
	FileHeader *multipart.FileHeader
}

func (f *File) Move(directory string, name ...string) (bool, error) {
	src, err := f.FileHeader.Open()
	if err != nil {
		return false, err
	}
	defer src.Close()

	fname := f.FileHeader.Filename

	if len(name) > 0 {
		fname = name[0]
	}

	dst := path.Join(directory, fname)

	out, err := os.Create(dst)
	if err != nil {
		return false, err
	}
	defer out.Close()

	io.Copy(out, src)

	return true, nil
}
