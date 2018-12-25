package route

import (
	"regexp"
	"strings"
)

// Rule Route rule
type Rule struct {
	Method         []string
	Pattern        string
	Handler        interface{}
	ParameterNames []string
	Parameters     map[string]string
	Compiled       *Compiled
}

type Compiled struct {
	Regex string
}

// Matches Determine if the rule matches given request.
func (r *Rule) Matches(method, path string) bool {
	r.compile()

	if false == matchMethods(method, r.Method) {
		return false
	}

	if false == matchPath(path, r.Compiled.Regex) {
		return false
	}

	return true
}

// Bind Bind the route to a given request for execution.
func (r *Rule) Bind(path string) {
	r.compile()

	path = "/" + strings.TrimLeft(path, "/")
	reg := regexp.MustCompile(r.Compiled.Regex)
	matches := reg.FindStringSubmatch(path)[1:]

	parameterNames := r.GetParameterNames()

	if len(parameterNames) == 0 {
		return
	}

	parameters := make(map[string]string)

	for k, v := range parameterNames {
		parameters[v] = matches[k]
	}

	r.Parameters = parameters
}

// GetParameterNames Get all of the parameter names for the rule.
func (r *Rule) GetParameterNames() []string {
	if r.ParameterNames != nil {
		return r.ParameterNames
	}
	r.ParameterNames = r.compileParameterNames()

	return r.ParameterNames
}

func (r *Rule) compile() {
	if r.Compiled != nil {
		return
	}

	r.Pattern = strings.Replace(r.Pattern, "/*", "/.*", -1)

	reg, _ := regexp.Compile(`\{\w+\}`)
	regex := reg.ReplaceAllString(r.Pattern, "([^/]+)")

	r.Compiled = &Compiled{
		Regex: "^" + regex + "$",
	}
}

func (r *Rule) compileParameterNames() []string {
	reg := regexp.MustCompile(`\{(.*?)\}`)
	matches := reg.FindAllStringSubmatch(r.Pattern, -1)

	var result []string
	for _, v := range matches {
		result = append(result, v[1])
	}

	return result
}
