package router

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"
)

var verbs = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

type Request interface {
	GetMethod() string
	GetPath() string
}

type Route struct {
	inited bool

	method      []string
	prefix      string
	pattern     string
	handler     interface{}
	middlewares []Middleware
	group       *Route

	collects []*Route

	rules    map[string]map[string]*Rule
	allRules map[string]*Rule
	// Current    *Rule
}

// New Create a new Route instance.
func New() *Route {
	route := &Route{
		rules: make(map[string]map[string]*Rule),
	}
	return route
}

// Dispatch Dispatch the request
func (r *Route) Dispatch(request Request) (*Rule, error) {
	rule, err := r.Match(request)
	if err != nil {
		return nil, err
	}

	rule.Bind(request.GetPath())

	return rule, nil
}

// Match Find the first rule matching a given request.
func (r *Route) Match(request Request) (*Rule, error) {
	for _, rule := range r.rules[request.GetMethod()] {
		if true == rule.Matches(request.GetMethod(), request.GetPath()) {
			return rule, nil
		}
	}
	return nil, errors.New("Not Found")
}

// AddRule Add a Rule to the Router.Rules
func (r *Route) AddRule(rule *Rule) *Rule {
	domainAndUri := rule.pattern
	method := ""
	for _, method = range rule.method {

		if _, ok := r.rules[method]; !ok {
			r.rules[method] = map[string]*Rule{
				domainAndUri: rule,
			}
		} else {
			r.rules[method][domainAndUri] = rule
		}

	}

	if r.allRules == nil {
		r.allRules = map[string]*Rule{}
	}
	r.allRules[method+domainAndUri] = rule

	return rule
}

// Add Add a router
func (r *Route) Add(method []string, pattern string, handler interface{}) *Route {
	route := r.initRoute()
	route.method = method
	route.pattern = r.getPrefix(pattern)
	route.handler = handler
	return route
}

// Get Register a new GET rule with the router.
func (r *Route) Get(pattern string, handler interface{}) *Route {
	return r.Add(Method("GET", "HEAD"), pattern, handler)
}

// Head Register a new Head rule with the router.
func (r *Route) Head(pattern string, handler interface{}) *Route {
	return r.Add(Method("HEAD"), pattern, handler)
}

// Post Register a new POST rule with the router.
func (r *Route) Post(pattern string, handler interface{}) *Route {
	return r.Add(Method("POST"), pattern, handler)
}

// Put Register a new PUT rule with the router.
func (r *Route) Put(pattern string, handler interface{}) *Route {
	return r.Add(Method("PUT"), pattern, handler)
}

// Patch Register a new PATCH rule with the router.
func (r *Route) Patch(pattern string, handler interface{}) *Route {
	return r.Add(Method("PATCH"), pattern, handler)
}

// Delete Register a new DELETE rule with the router.
func (r *Route) Delete(pattern string, handler interface{}) *Route {
	return r.Add(Method("DELETE"), pattern, handler)
}

// Options Register a new OPTIONS rule with the router.
func (r *Route) Options(pattern string, handler interface{}) *Route {
	return r.Add(Method("OPTIONS"), pattern, handler)
}

// Any Register a new rule responding to all verbs.
func (r *Route) Any(pattern string, handler interface{}) *Route {
	return r.Add(verbs, pattern, handler)
}

// Static Register a new Static rule.
func (r *Route) Static(path, root string) {
	path = "/" + strings.Trim(path, "/") + "/*"

	h := NewStaticHandle(root)

	r.Get(path, h)
	r.Head(path, h)
}

// Statics Bulk register Static rule.
func (r *Route) Statics(statics map[string]string) {
	for path, root := range statics {
		r.Static(path, root)
	}
}

// Prefix Add a prefix to the route URI.
func (r *Route) Prefix(prefix string) *Route {
	route := r.initRoute()
	route.prefix = route.getPrefix(prefix)
	return route
}

// Group Create a route group
func (r *Route) Group(callback func(group *Route)) *Route {
	route := r.initRoute()
	group := route.cloneRoute()
	group.Middleware(route.middlewares...)
	callback(group)
	return route
}

// Middleware Set the middleware attached to the route.
func (r *Route) Middleware(middlewares ...Middleware) *Route {
	route := r.initRoute()
	for _, m := range middlewares {
		route.middlewares = append(route.middlewares, m)
	}
	return route
}

// Register Register route from the collect.
func (r *Route) Register() {
	r.register(r)
}

func (r *Route) Dump() []byte {
	var b bytes.Buffer
	for _, rule := range r.allRules {
		fmt.Fprintf(&b, "%s %s %T \r\n", strings.Join(rule.method, "|"), rule.pattern, rule.handler)
	}

	return b.Bytes()
}

func (r *Route) register(root *Route) {
	for _, route := range r.collects {
		route.prefix = r.getPrefix(route.prefix)

		var middlewares []Middleware
		for _, m := range r.middlewares {
			middlewares = append(middlewares, m)
		}
		for _, m := range route.middlewares {
			middlewares = append(middlewares, m)
		}
		route.middlewares = middlewares

		route.register(root)

		if route.handler == nil {
			continue
		}
		rule := &Rule{
			method:      route.method,
			pattern:     route.getPrefix(route.pattern),
			handler:     route.handler,
			middlewares: route.middlewares,
		}

		root.AddRule(rule)
	}
	r.collects = r.collects[0:0]
}

// initRoute Initialize a new Route if not initialized
func (r *Route) initRoute() *Route {
	route := r
	if !r.inited {
		route = &Route{
			inited: true,
			rules:  make(map[string]map[string]*Rule),
			// prefix:      r.prefix,
			// middlewares: r.middlewares,
		}
		r.collects = append(r.collects, route)
	}
	return route
}

func (r *Route) cloneRoute() *Route {
	route := &Route{
		inited: false,
		rules:  make(map[string]map[string]*Rule),
		// prefix:      r.prefix,
		// middlewares: r.middlewares,
	}
	r.collects = append(r.collects, route)

	return route
}

func (r *Route) getPrefix(pattern string) string {
	return path.Join("/", r.prefix, pattern)
}
