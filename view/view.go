package view

import (
	"bytes"
	"html/template"
	"path"
)

// the default View.
var view = &View{}

type View struct {
	tmpl *template.Template
}

func New() *View {
	v := &View{}
	return v
}

// ParseGlob creates a new Template and parses the template definitions from the
// files identified by the pattern, which must match at least one file.
func (v *View) ParseGlob(pattern string) {
	v.tmpl = template.Must(template.ParseGlob(pattern))
}

func (v *View) Render(name string, data interface{}) template.HTML {
	tmpl := v.tmpl
	if tmpl == nil {
		pattern := path.Join(path.Dir(name), "*")
		name = path.Base(name)
		tmpl = template.Must(template.ParseGlob(pattern))
	}
	var buf bytes.Buffer
	tmpl.ExecuteTemplate(&buf, name, data)

	return template.HTML(buf.String())
}

func Render(name string, data interface{}) template.HTML {
	return view.Render(name, data)
}

func ParseGlob(pattern string) {
	view.tmpl = template.Must(template.ParseGlob(pattern))
}
