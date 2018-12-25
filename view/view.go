package view

import (
	"bytes"
	"html/template"
	"path"
	"strings"

	"github.com/thinkoner/thinkgo/filesystem"
	"github.com/thinkoner/thinkgo/helper"
)

type View struct {
	path string
	// Engine
}

func NewView() *View {
	v := &View{}
	return v
}

func (v *View) SetPath(path string) {
	path = helper.WorkPath(path)

	v.path = strings.TrimRight(path, "/") + "/"
}

func (v *View) loadEngine() {

}

func (v *View) Render(name string, data interface{}) []byte {
	file := name

	if e, _ := filesystem.Exists(file); !e {
		file = path.Join(v.path, name)
	}

	t := template.New("")
	if _, err := t.ParseGlob(path.Join(v.path, "*")); err != nil {
		panic(err)
	}

	tmpl, err := t.ParseFiles(file)
	if err != nil {
		panic(err)
		// return context.ErrorResponse()
	}

	var buf bytes.Buffer

	tmpl.ExecuteTemplate(&buf, name, data)

	return buf.Bytes()
}
