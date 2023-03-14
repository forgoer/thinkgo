package thinkgo

import (
	"crypto/tls"
	"fmt"
	"github.com/forgoer/thinkgo/ctx"
	"github.com/forgoer/thinkgo/router"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/forgoer/thinkgo/think"
	"github.com/stretchr/testify/assert"
)

func testRequest(t *testing.T, method, url string, data url.Values, res *ctx.Response) {
	var err error
	var resp *http.Response

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	var body io.Reader
	if strings.ToUpper(method) == "GET" {

	} else {

	}

	contentType := "application/x-www-form-urlencoded"

	method = strings.ToUpper(method)
	switch method {
	case "GET":
		url = strings.TrimRight(url, "?") + "?" + data.Encode()
	case "POST", "PUT", "DELETE":
		if data != nil {
			body = strings.NewReader(data.Encode())
		}
	}

	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", contentType)
	resp, err = client.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	content, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)

	assert.Equal(t, res.GetCode(), resp.StatusCode)
	assert.Equal(t, res.GetContent(), string(content))
}

func TestRunWithPort(t *testing.T) {
	th := New()

	go func() {
		th.Route(func(route *router.Route) {
			route.Get("/", func(req *ctx.Request) *ctx.Response {
				return think.Text("it worked")
			})
		})
		// listen and serve on 0.0.0.0:9011
		th.Run(":9012")
	}()

	time.Sleep(1 * time.Second)

	testRequest(t, "get", "http://localhost:9012/", nil, ctx.NewResponse().SetContent("it worked"))
}

func TestThink_Run(t *testing.T) {
	th := New()

	go func() {
		th.Route(func(route *router.Route) {
			route.Get("/", func(req *ctx.Request) interface{} {
				return "it worked"
			})
			route.Get("/user/{name}", func(req *ctx.Request, name string) interface{} {
				return fmt.Sprintf("Hello %s !", name)
			})
			route.Post("/user", func(req *ctx.Request) interface{} {
				name, err := req.Post("name")
				if err != nil {
					panic(err)
				}
				return fmt.Sprintf("Create %s", name)
			})
			route.Delete("/user/{name}", func(name string) interface{} {
				return fmt.Sprintf("Delete %s", name)
			})
		})
		// listen and serve on 0.0.0.0:9011
		th.Run()
	}()

	time.Sleep(5 * time.Millisecond)

	testRequest(t, "get", "http://localhost:9011/", nil, ctx.NewResponse().SetContent("it worked"))
	testRequest(t, "get", "http://localhost:9011/user/thinkgo", nil, ctx.NewResponse().SetContent(fmt.Sprintf("Hello %s !", "thinkgo")))
	testRequest(t, "post", "http://localhost:9011/user", url.Values{"name": {"thinkgo"}}, ctx.NewResponse().SetContent(fmt.Sprintf("Create %s", "thinkgo")))
	testRequest(t, "delete", "http://localhost:9011/user/thinkgo", url.Values{"name": {"thinkgo"}}, ctx.NewResponse().SetContent(fmt.Sprintf("Delete %s", "thinkgo")))
}
