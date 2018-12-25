<h1 align="center">
  ThinkGo
</h1>

<p align="center">
  <strong>ThinkGo is a lightweight MVC framework written in Go (Golang).</strong>
</p>

## Installation

[![Join the chat at https://gitter.im/think-go/community](https://badges.gitter.im/think-go/community.svg)](https://gitter.im/think-go/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

The only requirement is the [Go Programming Language](https://golang.org/dl/)

```
go get -u github.com/thinkoner/thinkgo
```

## Quick start

```go
package main

import (
	"github.com/thinkoner/thinkgo"
	"fmt"
	"github.com/thinkoner/thinkgo/route"
	"github.com/thinkoner/thinkgo/context"
)

func main() {
	app := thinkgo.BootStrap()
	app.RegisterRoute(func(route *route.Route) {

		route.Get("/", func(req *context.Request) thinkgo.Response {
			return thinkgo.Text("Hello ThinkGo !")
		})

		route.Get("/ping", func(req *context.Request) thinkgo.Response {
			return thinkgo.Json(map[string]string{
				"message": "pong",
			})
		})

		// Dependency injection
		route.Get("/user/{name}", func(req *context.Request, name string) thinkgo.Response {
			return thinkgo.Text(fmt.Sprintf("Hello %s !", name))
		})
	})
	// listen and serve on 0.0.0.0:9011
	app.Run()
}
```

## Table of contents

- [Routing](#routing)
- [Controller](#controller)
- [Request](#http-request)
- [Response](#http-response)
- [View](#view)
- [HTTP Session](#http-session)

## Routing

#### Basic Routing

The most basic routes accept a URI and a Closure, providing a very simple and expressive method of defining routes:

```go
app.RegisterRoute(func(route *route.Route) {
    route.Get("/foo", func(req *context.Request) thinkgo.Response {
        return thinkgo.Text("Hello ThinkGo !")
    })
})
```

#### Available Router Methods

The router allows you to register routes that respond to any HTTP verb:

```go
route.Get("/someGet", getting)
route.Post("/somePost", posting)
route.Put("/somePut", putting)
route.Delete("/someDelete", deleting)
route.Patch("/somePatch", patching)
route.Options("/someOptions", options)
```

Sometimes you may need to register a route that responds to multiple HTTP verbs. You may even register a route that responds to all HTTP verbs using the `Any` method:

```go
route.Any("/someAny", any)
```

#### Parameters in path

Of course, sometimes you will need to capture segments of the URI within your route. For example, you may need to capture a user's ID from the URL. You may do so by defining route parameters:

```go
route.Get("/user/{id}", func(req *context.Request, id string) thinkgo.Response {
	return thinkgo.Text(fmt.Sprintf("User %s", id))
})
```

You may define as many route parameters as required by your route:

```go
route.Get("/posts/{post}/comments/{comment}", func(req *context.Request, postId, commentId string) thinkgo.Response {
	//
})
```

## Controller

#### Basic Controller

Below is an example of a basic controller class.

```go
package controller

import (
	"github.com/thinkoner/thinkgo"
	"github.com/thinkoner/thinkgo/context"
)

func Index(req *context.Request) thinkgo.Response {
	return thinkgo.Text("Hello ThinkGo !")
}

```

You can define a route to this controller like so:

```go
route.Get("/", controller.Index)
```

#### Resource Controller

This feature will be supported in a future release.

## HTTP Request

#### Accessing The Request

To obtain an instance of the current HTTP request via dependency injection

```go
func Handler(req *context.Request) thinkgo.Response {
	name := req.Input("name")
}
```

#### Dependency Injection & Route Parameters

If your controller method is also expecting input from a route parameter you should list your route parameters after the request dependencies. For example, you can access your route parameter `name` like so:

```go
route.Put("/user/{name}", func(req *context.Request, name string) thinkgo.Response {
	//
})
```

#### Request Path & Method

The `path` method returns the request's path information. So, if the incoming request is targeted at `http://domain.com/foo/bar`, the `path` method will return `foo/bar`:

```go
uri := req.GetPath()
```

The `method` method will return the HTTP verb for the request. 

```go
method := req.GetMethod();
```

### Retrieving Cookies From Requests

```go
name, _ := request.Cookie("name")
```

## HTTP Response

an HTTP Response Must implement the `thinkgo.Response` interface

### Creating Responses

a simple strings or json Response:

```go
thinkgo.Text("Hello ThinkGo !")

thinkgo.Json(map[string]string{
				"message": "pong",
			})
```

### Attaching Cookies To Responses

```go
response.Cookie("name", "alice")
```

## View

views are stored in the `views` directory, A simple view might look something like this:

`views/layout.html` like this:

```html
{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
</head>
<body>
    {{ template "content" .}}
</body>
</html>
{{ end }}
```

`views/tpl.html` like this:

```html
{{ define "content" }}
<h2>{{ .Message }}</h2>
{{ end }}
{{ template "layout" . }}
```

we may return it using the `Render` function like so:

```go
route.Get("/tpl", func(request *context.Request) thinkgo.Response {
	data := map[string]interface{}{"Title": "ThinkGo", "Message": "Hello ThinkGo !"}
	return thinkgo.Render("tpl.html", data)
})
```

## HTTP Session

retrieving Data like this:

```go
request.Session().Get("user")
```

storing Data like this:

```go
request.Session().Set("user", "alice")
```

## License

This project is licensed under the [Apache 2.0 license](LICENSE).

## Contact

If you have any issues or feature requests, please contact us. PR is welcomed.
- https://github.com/thinkoner/thinkgo/issues
- duanpier@gmail.com
