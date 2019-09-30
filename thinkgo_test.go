package thinkgo

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thinkoner/thinkgo/context"
	"github.com/thinkoner/thinkgo/router"
)

func testRequest(t *testing.T, url string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)
	assert.Equal(t, "it worked", string(body), "resp body should match")
	assert.Equal(t, "200 OK", resp.Status, "should get a 200")
}

func TestRunEmpty(t *testing.T) {
	think := BootStrap()

	go func() {
		think.RegisterRoute(func(route *router.Route) {
			route.Get("/", func(req *context.Request) *context.Response {
				return Text("it worked")
			})
		})
		// listen and serve on 0.0.0.0:9011
		think.Run()
	}()

	time.Sleep(5 * time.Millisecond)

	testRequest(t, "http://localhost:9011/")
}

func TestRunWithPort(t *testing.T) {
	think := BootStrap()

	go func() {
		think.RegisterRoute(func(route *router.Route) {
			route.Get("/", func(req *context.Request) *context.Response {
				return Text("it worked")
			})
		})
		// listen and serve on 0.0.0.0:9011
		think.Run(":8100")
	}()

	time.Sleep(5 * time.Millisecond)

	testRequest(t, "http://localhost:8100/")
}
