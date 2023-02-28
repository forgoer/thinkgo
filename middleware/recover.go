package middleware

import (
	"fmt"
	"net/http/httputil"
	"runtime"
	"strings"

	. "github.com/forgoer/thinkgo/contracts"
	"github.com/forgoer/thinkgo/ctx"
)

type Recover struct {
	app Application
}

// New The default Recover handler
func (h *Recover) New(app Application) {
	h.app = app
}

// Handle an incoming request.
func (h *Recover) Handle(req *ctx.Request, next Next) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}

				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}

			httpRequest, _ := httputil.DumpRequest(req.Request, false)

			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					headers[idx] = current[0] + ": *"
				}
			}

			logMessage := fmt.Sprintf("Recovered at Request: %s\n", strings.Join(headers, "\r\n"))
			logMessage += fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)

			response := ctx.ErrorResponse()
			if h.app.Debug() {
				response.SetContent(logMessage)
			}
			result = response
		}
	}()

	return next(req)
}
