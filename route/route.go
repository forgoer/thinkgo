package route

import (
	"errors"
	"strings"
)

var verbs = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

type Request interface {
	GetMethod() string
	GetPath() string
}

type Route struct {
	Rules   map[string][]*Rule
	Current *Rule
}

func NewRoute() *Route {
	route := &Route{
		Rules: make(map[string][]*Rule),
	}
	return route
}

// Dispatch Dispatch the request to the application.
func (r *Route) Dispatch(request Request) (*Rule, error) {
	rule, err := r.Match(request)

	if err != nil {
		return nil, err
	}

	rule.Bind(request.GetPath())

	r.Current = rule

	return rule, nil
}

// Match Find the first rule matching a given request.
func (r *Route) Match(request Request) (*Rule, error) {
	for _, rule := range r.Rules[request.GetMethod()] {
		if true == rule.Matches(request.GetMethod(), request.GetPath()) {
			return rule, nil
		}
	}
	return nil, errors.New("Not Found")
}

// Add Add a route to the Router.Rules
func (r *Route) Add(method []string, pattern string, handler interface{}) {
	rule := &Rule{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	}
	var rules []*Rule

	for _, m := range rule.Method {
		if _, ok := r.Rules[m]; ok {
			rules = r.Rules[m]
		}
		r.Rules[m] = append(rules, rule)
	}
}

//Get Register a new GET rule with the router.
func (r *Route) Get(pattern string, handler interface{}) {
	r.Add(Method("GET"), pattern, handler)
}

//Head Register a new Head rule with the router.
func (r *Route) Head(pattern string, handler interface{}) {
	r.Add(Method("HEAD"), pattern, handler)
}

//Post Register a new POST rule with the router.
func (r *Route) Post(pattern string, handler interface{}) {
	r.Add(Method("POST"), pattern, handler)
}

//Put Register a new PUT rule with the router.
func (r *Route) Put(pattern string, handler interface{}) {
	r.Add(Method("PUT"), pattern, handler)
}

//Patch Register a new PATCH rule with the router.
func (r *Route) Patch(pattern string, handler interface{}) {
	r.Add(Method("PATCH"), pattern, handler)
}

//Delete Register a new DELETE rule with the router.
func (r *Route) Delete(pattern string, handler interface{}) {
	r.Add(Method("DELETE"), pattern, handler)
}

//Options Register a new OPTIONS rule with the router.
func (r *Route) Options(pattern string, handler interface{}) {
	r.Add(Method("OPTIONS"), pattern, handler)
}

func (r *Route) Static(path, root string) {
	path = "/" + strings.Trim(path, "/") + "/*"

	h := NewStaticHandle(root)

	r.Get(path, h)
	r.Head(path, h)
}

func (r *Route) Statics(statics map[string]string) {
	for path, root := range statics {
		r.Static(path, root)
	}
}

//Any Register a new rule responding to all verbs.
func (r *Route) Any(pattern string, handler interface{}) {
	r.Add(verbs, pattern, handler)
}
