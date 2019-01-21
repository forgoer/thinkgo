package thinkgo

import (
	"encoding/json"

	"github.com/thinkoner/thinkgo/context"
)

// Json Create a new HTTP Response with JSON data
func Json(v interface{}) *context.Response {
	r := context.NewResponse()
	r.SetContentType("application/json")
	c, _ := json.Marshal(v)
	r.SetContent(string(c))
	return r
}

// Text Create a new HTTP Response with TEXT data
func Text(s string) *context.Response {
	r := context.NewResponse()
	r.SetContentType("text/plain")
	r.SetContent(s)
	return r
}

// Text Create a new HTTP Response with HTML data
func Html(s string) *context.Response {
	r := context.NewResponse()
	r.SetContent(s)
	return r
}

// Render Create a new HTTP Response with the template
func Render(name string, data interface{}) *context.Response {
	b := application.GetView().Render(name, data)
	return Html(string(b))
}
