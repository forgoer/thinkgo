package context

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//Request HTTP request
type Request struct {
	Request       *http.Request
	method        string
	path          string
	query         map[string]string
	post          map[string]string
	files         map[string]*File
	session       Session
	CookieHandler *Cookie
}

// NewRequest create a new HTTP request from *http.Request
func NewRequest(req *http.Request) *Request {

	return &Request{
		Request: req,
		method:  req.Method,
		path:    req.URL.Path,
		query:   parseQuery(req.URL.Query()),
		post:    parsePost(req),
	}
}

// GetMethod get the request method.
func (r *Request) GetMethod() string {
	return r.method
}

// GetPath get the request path.
func (r *Request) GetPath() string {
	return r.path
}

// GetHttpRequest get Current *http.Request
func (r *Request) GetHttpRequest() *http.Request {
	return r.Request
}

//IsMethod checks if the request method is of specified type.
func (r *Request) IsMethod(m string) bool {
	return strings.ToUpper(m) == r.GetMethod()
}

// Query returns a query string item from the request.
func (r *Request) Query(key string, value ...string) (string, error) {
	if v, ok := r.query[key]; ok {
		return v, nil
	}
	if len(value) > 0 {
		return value[0], nil
	}
	return "", errors.New("named query not present")
}

// Input returns a input item from the request.
func (r *Request) Input(key string, value ...string) (string, error) {
	if v, ok := r.post[key]; ok {
		return v, nil
	}

	if v, ok := r.query[key]; ok {
		return v, nil
	}

	if len(value) > 0 {
		return value[0], nil
	}
	return "", errors.New("named input not present")
}

// Input returns a post item from the request.
func (r *Request) Post(key string, value ...string) (string, error) {
	if v, ok := r.post[key]; ok {
		return v, nil
	}

	if len(value) > 0 {
		return value[0], nil
	}
	return "", errors.New("named post not present")
}

//Cookie Retrieve a cookie from the request.
func (r *Request) Cookie(key string, value ...string) (string, error) {
	var err error
	key = r.CookieHandler.Config.Prefix + key
	cookie, err := r.Request.Cookie(key)
	if err == nil {
		c, _ := url.QueryUnescape(cookie.Value)
		return c, err
	}
	if len(value) > 0 {
		return value[0], nil
	}
	return "", err
}

// File returns a file from the request.
func (r *Request) File(key string) (*File, error) {
	if f, ok := r.files[key]; ok {
		return f, nil
	}

	_, fh, err := r.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	r.files[key] = &File{fh}

	return r.files[key], nil
}

// AllFiles returns all files from the request.
func (r *Request) AllFiles() (map[string]*File, error) {
	err := r.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}
	if r.Request.MultipartForm != nil || r.Request.MultipartForm.File != nil {
		for key, fh := range r.Request.MultipartForm.File {
			r.files[key] = &File{fh[0]}
		}
	}
	return r.files, nil
}

// All get all of the input and query for the request.
func (r *Request) All(keys ...string) map[string]string {
	all := mergeForm(r.query, r.post)

	if len(keys) == 0 {
		return all
	}

	result := make(map[string]string)

	for _, key := range keys {
		if v, ok := all[key]; ok {
			result[key] = v
		} else {
			result[key] = ""
		}
	}

	return result
}

//Only get a subset of the items from the input data.
func (r *Request) Only(keys ...string) map[string]string {
	all := r.All()

	result := make(map[string]string)

	for _, key := range keys {
		if v, ok := all[key]; ok {
			result[key] = v
		}
	}

	return result
}

//Except Get all of the input except for a specified array of items.
func (r *Request) Except(keys ...string) map[string]string {
	all := r.All()

	for _, key := range keys {
		if _, ok := all[key]; ok {
			delete(all, key)
		}
	}

	return all
}

// Has Determine if the request contains a given input item key.
func (r *Request) Exists(keys ...string) bool {
	all := r.All()

	for _, key := range keys {
		if _, ok := all[key]; !ok {
			return false
		}
	}

	return true
}

// Filled Determine if the request contains a non-empty value for an input item.
func (r *Request) Has(keys ...string) bool {
	all := r.All()

	for _, key := range keys {
		if _, ok := all[key]; !ok {
			return false
		}
		if len(all[key]) == 0 {
			return false
		}
	}

	return true
}

//Url get the URL (no query string) for the request.
func (r *Request) Url() string {
	return r.Request.URL.Path
}

// FullUrl get the full URL for the request.
func (r *Request) FullUrl() string {
	return r.Url() + "?" + r.Request.URL.RawQuery
}

// Path get the current path info for the request.
func (r *Request) Path() string {
	return r.path
}

// Method get the current method for the request.
func (r *Request) Method() string {
	return r.method
}

// GetContent Returns the request body content.
func (r *Request) GetContent() ([]byte, error) {
	var body []byte

	if body == nil {
		body, err := ioutil.ReadAll(r.Request.Body)
		if err != nil {
			return nil, err
		}
		r.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	return body, nil
}

// Session get the session associated with the request.
func (r *Request) Session() Session {
	return r.session
}

// Session set the session associated with the request.
func (r *Request) SetSession(s Session) {
	r.session = s
}

func parseQuery(q url.Values) map[string]string {
	query := make(map[string]string)
	for k, v := range q {
		query[k] = v[0]
	}
	return query
}

func parsePost(r *http.Request) map[string]string {
	post := make(map[string]string)

	r.ParseForm()
	for k, v := range r.PostForm {
		post[k] = v[0]
	}

	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm != nil {
		for k, v := range r.MultipartForm.Value {
			post[k] = v[0]
		}
	}

	return post
}

func mergeForm(slices ...map[string]string) map[string]string {
	r := make(map[string]string)

	for _, slice := range slices {
		for k, v := range slice {
			r[k] = v
		}
	}

	return r
}
