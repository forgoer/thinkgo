package thinkgo

import (
	"encoding/json"
	"net/http"

	"github.com/thinkoner/thinkgo/context"
)

// Response an HTTP response interface
type Response interface {
	Send(w http.ResponseWriter)
}

// NewResponse Create a new HTTP Response
func NewResponse() *context.Response {
	r := &context.Response{}
	r.SetContentType("text/html")
	r.SetCharset("utf-8")
	r.SetCode(200)
	r.CookieHandler = parseCookieHandler()
	return r
}

// NotFoundResponse Create a new HTTP NotFoundResponse
func NotFoundResponse() Response {
	r := NewResponse()
	r.SetContent("Not Found")
	r.SetCode(404)
	return r
}

// NotFoundResponse Create a new HTTP ErrorResponse
func ErrorResponse() Response {
	r := NewResponse()
	r.SetContent("Server Error")
	r.SetCode(500)
	return r
}

// Json Create a new HTTP Response with JSON data
func Json(v interface{}) *context.Response {
	r := NewResponse()
	r.SetContentType("application/json")
	c, _ := json.Marshal(v)
	r.SetContent(string(c))
	return r
}

// Text Create a new HTTP Response with TEXT data
func Text(s string) *context.Response {
	r := NewResponse()
	r.SetContentType("text/plain")
	r.SetContent(s)
	return r
}

// Text Create a new HTTP Response with HTML data
func Html(s string) *context.Response {
	r := NewResponse()
	r.SetContent(s)
	return r
}

// Render Create a new HTTP Response with the template
func Render(name string, data interface{}) *context.Response {
	b := application.View.Render(name, data)
	return Html(string(b))
}
