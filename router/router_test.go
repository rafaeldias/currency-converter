package router

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"testing"
)

func TestSanatizePath(t *testing.T) {
	var testCase = []struct {
		path, want string
	}{
		{"", "/"},
		{"ping", "/ping"},
	}

	for _, tc := range testCase {
		p := sanatizePath(tc.path)

		if p != tc.want {
			t.Errorf("got: %s, want: %s", p, tc.want)
		}

	}
}

func TestParamsFromPath(t *testing.T) {
	var testCase = []struct {
		path string
		want []string
	}{
		{"/users/:id", []string{"id"}},
		{"/users/:id/files/:name", []string{"id", "name"}},
	}

	for _, tc := range testCase {
		p := paramsFromPath(tc.path)

		if !reflect.DeepEqual(p, tc.want) {
			t.Errorf("got: %s, want: %s", p, tc.want)
		}
	}
}

func TestPatternFromPath(t *testing.T) {
	var testCase = []struct {
		path string
		want *regexp.Regexp
	}{
		{"/users", regexp.MustCompile("^/users$")},
		{"/users/:id", regexp.MustCompile("^/users/([a-z0-9]+)$")},
		{"/users/:id/files/:name", regexp.MustCompile("^/users/([a-z0-9]+)/files/([a-z0-9]+)$")},
	}

	for _, tc := range testCase {
		r := patternFromPath(tc.path)

		if !reflect.DeepEqual(r, tc.want) {
			t.Errorf("got: %s, want: %s", r, tc.want)
		}
	}
}

func noopHandler(c HTTPContexter) {}

func TestGet(t *testing.T) {
	var testCases = []struct {
		path    string
		handler Handler
		want    Route
	}{
		{
			"/users/:id",
			noopHandler,
			Route{
				Handler: Handler(noopHandler),
				Method:  http.MethodGet,
				Params:  []string{"id"},
				Pattern: regexp.MustCompile("^/users/([a-z0-9]+)$"),
			},
		},
	}

	for _, tc := range testCases {
		r := New()

		r.Get(tc.path, tc.handler)

		if len(r.Routes) < 1 || !reflect.DeepEqual(r.Routes[0].Pattern, tc.want.Pattern) {
			t.Errorf("got: %v, want: %v", r.Routes, tc.want)
		}
	}
}

func TestUse(t *testing.T) {
	var testCases = []struct {
		handler Handler
	}{
		{noopHandler},
	}

	for _, tc := range testCases {
		r := New()

		r.Use(tc.handler)

		if len(r.Middlewares) < 1 {
			t.Errorf("got: %v, want: %v", r.Middlewares, tc.handler)
		}
	}
}

type testHandler struct {
	called bool
}

func (t *testHandler) Handle(c HTTPContexter) {
	t.called = true
}

func TestServeHTTP(t *testing.T) {
	var testCases = []struct {
		handler *testHandler
		pattern string
		path    string
	}{
		{&testHandler{}, "/users/:id", "/users/1"},
		{&testHandler{}, "/users/:id", "/users/xyz"},
	}

	for _, tc := range testCases {
		url := &url.URL{Path: tc.path}
		rw := &responseWriterTest{}
		req := &http.Request{Method: http.MethodGet, URL: url}

		r := New()

		r.Get(tc.pattern, tc.handler.Handle)

		r.ServeHTTP(rw, req)

		if !tc.handler.called {
			t.Errorf("got: %v, want: true", tc.handler.called)
		}
	}
}

func TestSetHTTPParams(t *testing.T) {
	var testCases = []struct {
		ctx     HTTPContexter
		matches []string
		r       Route
	}{
		{
			newHTTPContext(nil, nil),
			[]string{"", "1"}, // simulate FindStringSubmatch return
			Route{
				Params: []string{"id"},
			},
		},
	}

	for _, tc := range testCases {
		setHTTPParams(tc.ctx, tc.r, tc.matches)

		if p := tc.ctx.Param("id"); p != "1" {
			t.Errorf("got: %s, want: %s", p, tc.matches[1])
		}
	}
}

func TestExecMiddlewares(t *testing.T) {
	var testCases = []struct {
		h    *testHandler
		want bool
	}{
		{&testHandler{}, true},
	}

	for _, tc := range testCases {
		r := New()
		r.Use(tc.h.Handle)
		r.execMiddlewares(newHTTPContext(nil, nil))

		if tc.h.called != tc.want {
			t.Errorf("got: %v, want: %v", tc.h.called, tc.want)
		}
	}
}

func TestServeHTTPMiddleware(t *testing.T) {
	var testCases = []struct {
		handler *testHandler
		pattern string
		path    string
	}{
		{&testHandler{}, "/users/:id", "/users/1"},
	}

	for _, tc := range testCases {
		url := &url.URL{Path: tc.path}
		rw := &responseWriterTest{}
		req := &http.Request{Method: http.MethodGet, URL: url}

		r := New()

		r.Use(tc.handler.Handle)
		r.Get(tc.pattern, noopHandler)

		r.ServeHTTP(rw, req)

		if !tc.handler.called {
			t.Errorf("got: %v, want: true", tc.handler.called)
		}
	}
}
