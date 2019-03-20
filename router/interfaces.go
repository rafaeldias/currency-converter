package router

import "net/http"

// HTTPContexter is the context of the HTTP request
type HTTPContexter interface {
	JSONer
	Getter
	Setter

	// Header returns the header that will be sent in response.
	// Enabling the user to set custom headers to the response.
	Header() http.Header
	// Request return the *http.Request of the currenct request
	Request() *http.Request
	// Sets param sent via HTTP request URL
	SetParam(name, value string)
	// Returns param sent via HTTP request URL
	Param(name string) string
	// Write response
	Write(code int, b []byte) error
}

// Getter is an interface used for getting contextual data
type Getter interface {
	Get(s string) interface{}
}

// Setter is an interface used for setting contextual data
type Setter interface {
	Set(s string, v interface{})
}

// JSONer converts v into a slice of bytes in JSON format
// if jsonp is true, json will be enclosed by the jsonp
// paramter, if it exists in the http request.
type JSONer interface {
	JSON(v interface{}, jsonp bool) []byte
}
