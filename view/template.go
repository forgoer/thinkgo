package view

//import (
//	"strings"
//	"github.com/thinkoner/thinkgo/app"
//)
//
//type Template struct {
//	ViewPath string
//}
//
//func NewTemplate() *Template {
//	t := &Template{}
//	t.loadDefault()
//
//	return t
//}
//
//func (t *Template) loadDefault() {
//	path, err := app.Config.GetString("view.path")
//	if err != nil {
//		panic(err)
//	}
//
//	path = WorkPath(path)
//
//	t.ViewPath = strings.TrimRight(path, "/") + "/"
//}
