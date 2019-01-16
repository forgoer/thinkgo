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
}

// GetContentType sets the Content-Type on the response.
func (r *Response) SetContentType(val string) {
	r.contentType = val
}

// GetContentType sets the Charset on the response.
func (r *Response) SetCharset(val string) {
	r.charset = val
}

// SetCode sets the status code on the response.
func (r *Response) SetCode(val int) {
	r.code = val
}

// SetContent sets the content on the response.
func (r *Response) SetContent(val string) {
	r.content = val
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

	w.Header().Set("Content-Type", r.GetContentType()+";"+" charset="+r.GetCharset())

	w.WriteHeader(r.GetCode())

	w.Write([]byte(r.GetContent()))
}
