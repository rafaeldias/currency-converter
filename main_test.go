package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/rafaeldias/currency-converter/router"
)

type ctxHTTPMainTest struct {
	header http.Header
	data   map[string]interface{}
}

func (c *ctxHTTPMainTest) JSON(v interface{}, jsonp bool) []byte {
	return []byte{}
}

func (c *ctxHTTPMainTest) Get(s string) interface{} {
	return c.data[s]
}

func (c *ctxHTTPMainTest) Set(s string, v interface{}) {
	c.data[s] = v
}

func (c *ctxHTTPMainTest) Request() *http.Request {
	return nil
}

func (c *ctxHTTPMainTest) SetParam(name, value string) {}

func (c *ctxHTTPMainTest) Param(name string) string {
	return ""
}

func (c *ctxHTTPMainTest) Write(code int, b []byte) error {
	return nil
}

func (c *ctxHTTPMainTest) Header() http.Header {
	return c.header
}

func TestGetEnv(t *testing.T) {
	var testCases = []struct {
		env   string
		value string
		def   string
		want  string
	}{
		{"x", "y", "", "y"},
		{"x", "", "z", "z"},
		{"x", "y", "z", "y"},
	}

	for _, tc := range testCases {
		os.Setenv(tc.env, tc.value)
		env := getEnv(tc.env, tc.def)

		if env != tc.want {
			t.Errorf("got: %s, want: %s", env, tc.want)
		}
	}
}

func TestCorsHeader(t *testing.T) {
	var header = http.Header{}
	var ctx = &ctxHTTPMainTest{header: header}
	var handler = CORSHeaders()

	handler(ctx)

	if h := ctx.Header().Get("Access-Control-Allow-Origin"); h != "*" {
		t.Errorf("got: %s, want: *", h)
	}
}

func TestCurrencyLayerMiddleware(t *testing.T) {
	var data = map[string]interface{}{}
	var ctx = &ctxHTTPMainTest{data: data}

	var handler = currencyLayerMiddleware("x", "y")

	handler(ctx)

	if _, ok := data["currency"]; !ok {
		t.Errorf("got: %v, want: true", ok)
	}
}

func TestHandler(t *testing.T) {
	var h = handler("host", "key")

	if r, ok := h.(*router.Router); !ok {
		t.Errorf("got: %v, want: *router.Router", r)
	}

}
