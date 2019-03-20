package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	headerJSON  = "application/json"
	headerJSONP = "application/javascript"
)

// HTTPContext represents groups
type HTTPContext struct {
	http.ResponseWriter

	req *http.Request

	ctxData map[string]interface{}
	params  map[string]string
}

func newHTTPContext(rw http.ResponseWriter, req *http.Request) HTTPContexter {
	return &HTTPContext{rw, req, map[string]interface{}{}, map[string]string{}}
}

// JSON writes to ResponseWriter the v in JSON format with the code as StatusCode
func (c *HTTPContext) JSON(v interface{}, jsonp bool) []byte {
	b, _ := json.Marshal(v)

	c.Header().Set("Content-Type", headerJSON)

	if fn := c.req.URL.Query().Get("jsonp"); jsonp && fn != "" {
		c.Header().Set("Content-Type", headerJSONP)
		b = []byte(fmt.Sprintf("%s(%s)", fn, b))
	}

	return b
}

// Get returns contextual data with name s
func (c *HTTPContext) Get(s string) interface{} {
	v, ok := c.ctxData[s]
	if !ok {
		return nil
	}
	return v
}

// Set sets contextual data v with name s
func (c *HTTPContext) Set(s string, v interface{}) {
	c.ctxData[s] = v
}

// Request returns the current *http.Request
func (c *HTTPContext) Request() *http.Request {
	return c.req
}

// SetParam returns the parameter with the name n
func (c *HTTPContext) SetParam(name, value string) {
	c.params[name] = value

}

// Param returns the parameter with the name n
func (c *HTTPContext) Param(name string) string {
	p, ok := c.params[name]
	if !ok {
		return ""
	}
	return p
}

// Write writes uses b as the response to the request
func (c *HTTPContext) Write(code int, b []byte) error {
	c.WriteHeader(code)
	c.ResponseWriter.Write(b)

	return nil
}
