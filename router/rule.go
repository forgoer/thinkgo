package router

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/thinkoner/thinkgo/context"
)

// Rule Route rule
type Rule struct {
	middlewares    []Middleware
	method         []string
	pattern        string
	handler        interface{}
	parameterNames []string
	parameters     map[string]string
	Compiled       *Compiled
}

type Compiled struct {
	Regex string
}

// Matches Determine if the rule matches given request.
func (r *Rule) Matches(method, path string) bool {
	r.compile()

	if false == matchMethods(method, r.method) {
		return false
	}

	if false == matchPath(path, r.Compiled.Regex) {
		return false
	}

	return true
}

// Bind Bind the router to a given request for execution.
func (r *Rule) Bind(path string) {
	r.compile()

	path = "/" + strings.TrimLeft(path, "/")
	reg := regexp.MustCompile(r.Compiled.Regex)
	matches := reg.FindStringSubmatch(path)[1:]

	parameterNames := r.getParameterNames()

	if len(parameterNames) == 0 {
		return
	}

	parameters := make(map[string]string)

	for k, v := range parameterNames {
		parameters[v] = matches[k]
	}

	r.parameters = parameters
}

// Middleware Set the middleware attached to the rule.
func (r *Rule) Middleware(middlewares ...Middleware) *Rule {
	for _, m := range middlewares {
		r.middlewares = append(r.middlewares, m)
	}
	return r
}

// GatherRouteMiddleware Get all middleware, including the ones from the controller.
func (r *Rule) GatherRouteMiddleware() []Middleware {
	return r.middlewares
}

// Run Run the route action and return the response.
func (r *Rule) Run(request *context.Request) interface{} {
	handler := r.handler

	if handler != nil {
		t := reflect.TypeOf(handler)
		switch t.Kind() {
		case reflect.Func:
			v := reflect.ValueOf(handler).Call(
				parseParams(request, r.parameters),
			)
			handler = v[0].Interface()
		}
	}
	return handler
}

// getParameterNames Get all of the parameter names for the rule.
func (r *Rule) getParameterNames() []string {
	if r.parameterNames != nil {
		return r.parameterNames
	}
	r.parameterNames = r.compileParameterNames()

	return r.parameterNames
}

func (r *Rule) compile() {
	if r.Compiled != nil {
		return
	}

	r.pattern = strings.Replace(r.pattern, "/*", "/.*", -1)

	reg, _ := regexp.Compile(`\{\w+\}`)
	regex := reg.ReplaceAllString(r.pattern, "([^/]+)")

	r.Compiled = &Compiled{
		Regex: "^" + regex + "$",
	}
}

func (r *Rule) compileParameterNames() []string {
	reg := regexp.MustCompile(`\{(.*?)\}`)
	matches := reg.FindAllStringSubmatch(r.pattern, -1)

	var result []string
	for _, v := range matches {
		result = append(result, v[1])
	}

	return result
}

func (r *Rule) toString() string {
	return fmt.Sprint(r.method) + ": " + r.pattern
}

func parseParams(request *context.Request, parameters map[string]string) []reflect.Value {
	var params []reflect.Value
	params = append(params, reflect.ValueOf(request))

	for _, v := range parameters {
		params = append(params, reflect.ValueOf(v))
	}

	return params
}
