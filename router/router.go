package router

import (
	"github.com/forgoer/thinkgo/ctx"
)

// RunRoute Return the response for the given rule.
func RunRoute(request *ctx.Request, rule *Rule) interface{} {
	return PrepareResponse(
		request,
		rule,
		runMiddlewares(request, rule),
	)
}

// PrepareResponse Create a response instance from the given value.
func PrepareResponse(request *ctx.Request, rule *Rule, result interface{}) interface{} {
	return result
	//var response Response
	//switch result.(type) {
	//case Response:
	//	return result.(Response)
	//// case string:
	//case ctx.Handler:
	//	return response
	//default:
	//	response = ctx.NewResponse().SetContent(fmt.Sprint(result))
	//}
	//return response
}

// runMiddlewares Run the given route within Middlewares instance.
func runMiddlewares(request *ctx.Request, rule *Rule) interface{} {
	pipeline := NewPipeline()

	for _, m := range rule.GatherRouteMiddleware() {
		pipeline.Pipe(m)
	}
	return pipeline.Passable(request).Run(func(request *ctx.Request, next Closure) interface{} {
		return PrepareResponse(
			request,
			rule,
			rule.Run(request),
		)
	})
}
