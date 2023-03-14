package think

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/forgoer/thinkgo/ctx"
)

// Json Create a new HTTP Response with JSON data
func Json(v interface{}) *ctx.Response {
	c, _ := json.Marshal(v)
	return ctx.NewResponse().SetContent(string(c)).SetContentType("application/json")
}

// Text Create a new HTTP Response with TEXT data
func Text(s string) *ctx.Response {
	return ctx.NewResponse().SetContent(s).SetContentType("text/plain")
}

// Text Create a new HTTP Response with HTML data
func Html(s string) *ctx.Response {
	return ctx.NewResponse().SetContent(s)
}

func Response(v interface{}) *ctx.Response {
	r := ctx.NewResponse()
	r.SetContentType("text/plain")

	var content string
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Bool:
		content = fmt.Sprintf("%t", v)
	case reflect.String:
		content = fmt.Sprintf("%s", v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		content = fmt.Sprintf("%d", v)
	case reflect.Float32, reflect.Float64:
		content = fmt.Sprintf("%v", v)
	default:
		r.SetContentType("application/json")
		b, err := json.Marshal(v)
		if err != nil {
			panic(errors.New(`The Response content must be a string, numeric, boolean, slice, map or can be encoded as json a string, "` + t.Name() + `" given.`))
		}
		content = string(b)
	}
	r.SetContent(content)

	return r
}
