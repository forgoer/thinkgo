package router

import (
	"fmt"
	"net/http"

	"github.com/thinkoner/thinkgo/context"
)

// RunRoute Return the response for the given rule.
func RunRoute(request *context.Request, rule *Rule) interface{} {
	return PrepareResponse(
		request,
		rule,
		runMiddlewares(request, rule),
	)
}

// PrepareResponse Create a response instance from the given value.
func PrepareResponse(request *context.Request, rule *Rule, result interface{}) interface{} {
	var response Response
	switch result.(type) {
	case Response:
		return result.(Response)
	// case string:
	case http.Handler:
		return response
	default:
		// t := reflect.TypeOf(response)
		// switch t.Kind() {
		// case reflect.Func:
		// 	v := reflect.ValueOf(response).Call(
		// 		parseParams(request, rule.Parameters),
		// 	)
		// 	response = v[0].Interface()
		// default:
		//
		// }
		response = context.NewResponse().SetContent(fmt.Sprint(result))
	}
	return response
}

// runMiddlewares Run the given route within Middlewares instance.
func runMiddlewares(request *context.Request, rule *Rule) interface{} {
	pipeline := NewPipeline()

	for _, m := range rule.GatherRouteMiddleware() {
		pipeline.Pipe(m)
	}
	return pipeline.Passable(request).Run(func(request *context.Request, next Closure) interface{} {
		return PrepareResponse(
			request,
			rule,
			rule.Run(request),
		)
	})
}
