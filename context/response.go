package context

import (
	"net/http"
)

type Response struct {
	// Writer      context.ResponseWriter
	contentType   string
	charset       string
	code          int
	content       string
	cookies       map[string]*http.Cookie
	CookieHandler *Cookie
	Header        *http.Header
}

// GetContentType sets the Content-Type on the response.
func (r *Response) SetContentType(val string) *Response {
	r.contentType = val
	return r
}

// GetContentType sets the Charset on the response.
func (r *Response) SetCharset(val string) *Response {
	r.charset = val
	return r
}

// SetCode sets the status code on the response.
func (r *Response) SetCode(val int) *Response {
	r.code = val
	return r
}

// SetContent sets the content on the response.
func (r *Response) SetContent(val string) *Response {
	r.content = val
	return r
}

// GetContentType get the Content-Type on the response.
func (r *Response) GetContentType() string {
	return r.contentType
}

// GetContentType get the Charset on the response.
func (r *Response) GetCharset() string {
	return r.charset
}

// GetCode get the response status code.
func (r *Response) GetCode() int {
	return r.code
}

// GetCode get the response content.
func (r *Response) GetContent() string {
	return r.content
}

// Cookie Add a cookie to the response.
func (r *Response) Cookie(name interface{}, params ...interface{}) error {
	cookie, err := r.CookieHandler.Set(name, params...)

	if err != nil {
		if r.cookies == nil {
			r.cookies = make(map[string]*http.Cookie)
		}
		r.cookies[cookie.Name] = cookie
	}

	return err
}

// Send Sends HTTP headers and content.
func (r *Response) Send(w http.ResponseWriter) {
	for _, cookie := range r.cookies {
		http.SetCookie(w, cookie)
	}
	for key, value := range *r.Header {
		for _, val := range value {
			w.Header().Add(key, val)
		}
	}
	w.Header().Set("Content-Type", r.GetContentType()+";"+" charset="+r.GetCharset())
	// r.Header.Write(w)
	w.WriteHeader(r.GetCode())
	w.Write([]byte(r.GetContent()))
}

// NewResponse Create a new HTTP Response
func NewResponse() *Response {
	r := &Response{
		Header: &http.Header{},
	}
	r.SetContentType("text/html")
	r.SetCharset("utf-8")
	r.SetCode(http.StatusOK)
	r.CookieHandler = ParseCookieHandler()
	return r
}

// NotFoundResponse Create a new HTTP NotFoundResponse
func NotFoundResponse() *Response {
	r := NewResponse()
	r.SetContent("Not Found")
	r.SetCode(http.StatusNotFound)
	return r
}

// NotFoundResponse Create a new HTTP Error Response
func ErrorResponse() *Response {
	r := NewResponse()
	r.SetContent("Server Error")
	r.SetCode(http.StatusInternalServerError)
	return r
}

// Redirect Create a new HTTP Redirect Response
func Redirect(to string) *Response {
	r := NewResponse()
	r.Header.Set("Location", to)
	r.SetCode(http.StatusMovedPermanently)
	return r
}
