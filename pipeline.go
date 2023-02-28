package thinkgo

import (
	"container/list"
	"html/template"
	"net/http"

	"github.com/forgoer/thinkgo/ctx"

	"github.com/forgoer/thinkgo/contracts"
	"github.com/forgoer/thinkgo/router"
	"github.com/forgoer/thinkgo/think"
)

type Pipeline struct {
	handlers []contracts.Middleware
	pipeline *list.List
	passable *ctx.Request
}

// Pipeline returns a new Pipeline
func NewPipeline() *Pipeline {
	p := &Pipeline{
		pipeline: list.New(),
	}
	return p
}

// Pipe Push a Middleware Handler to the pipeline
func (p *Pipeline) Pipe(m contracts.Middleware) *Pipeline {
	p.pipeline.PushBack(m)
	return p
}

// Pipe Batch push Middleware Handlers to the pipeline
func (p *Pipeline) Through(hls []contracts.Middleware) *Pipeline {
	for _, hl := range hls {
		p.Pipe(hl)
	}
	return p
}

// Passable set the request being sent through the pipeline.
func (p *Pipeline) Passable(passable *ctx.Request) *Pipeline {
	p.passable = passable
	return p
}

// Then Run the pipeline with a final destination callback.
func (p *Pipeline) Then(hl contracts.Middleware) {
	p.Pipe(hl)
}

// Run the pipeline
func (p *Pipeline) Run() interface{} {
	var result interface{}
	e := p.pipeline.Front()
	if e != nil {
		result = p.handler(p.passable, e)
	}
	return result
}

// ServeHTTP
func (p *Pipeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := ctx.NewRequest(r)
	request.CookieHandler = ctx.ParseCookieHandler()
	p.Passable(request)

	result := p.Run()

	switch result := result.(type) {
	case router.Response:
		result.Send(w)
	case template.HTML:
		think.Html(string(result)).Send(w)
	case http.Handler:
		result.ServeHTTP(w, r)
	default:
		think.Response(result).Send(w)
	}
}

func (p *Pipeline) handler(passable *ctx.Request, e *list.Element) interface{} {
	if e == nil {
		return nil
	}
	hl := e.Value.(contracts.Middleware)
	result := hl.Handle(passable, p.closure(e))
	return result
}

func (p *Pipeline) closure(e *list.Element) contracts.Next {
	return func(req *ctx.Request) interface{} {
		e = e.Next()
		return p.handler(req, e)
	}
}
