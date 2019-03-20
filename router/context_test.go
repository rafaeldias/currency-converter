package router

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type responseWriterTest struct {
	writtenBytes []byte
	statusCode   int
	h            http.Header
}

func (rw *responseWriterTest) Header() http.Header {
	return rw.h
}

func (rw *responseWriterTest) Write(b []byte) (int, error) {
	rw.writtenBytes = b
	return len(b), nil
}

func (rw *responseWriterTest) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	return
}

func TestJSON(t *testing.T) {
	var testCases = []struct {
		json interface{}
		want []byte
	}{
		{struct {
			Test int `json:"test"`
		}{1}, []byte("{\"test\":1}")},
	}

	for _, tc := range testCases {
		var c = newHTTPContext(&responseWriterTest{h: http.Header{}}, &http.Request{URL: &url.URL{}})
		var b = c.JSON(tc.json, false)

		if !reflect.DeepEqual(b, tc.want) {
			t.Errorf("got: %s, want: %s", b, tc.want)
		}

	}
}

func TestJSONP(t *testing.T) {
	var testCases = []struct {
		json     interface{}
		jsonp    bool
		jsonpArg string
		want     []byte
	}{
		{struct {
			Test int `json:"test"`
		}{1}, true, "testing", []byte("testing({\"test\":1})")},
		{struct {
			Test int `json:"test"`
		}{1}, true, "", []byte("{\"test\":1}")},
	}

	for _, tc := range testCases {
		var values = url.Values{
			"jsonp": {tc.jsonpArg},
		}
		var URL = &url.URL{
			RawQuery: values.Encode(),
		}
		var c = newHTTPContext(&responseWriterTest{h: http.Header{}}, &http.Request{URL: URL})
		var b = c.JSON(tc.json, tc.jsonp)

		if !reflect.DeepEqual(b, tc.want) {
			t.Errorf("got: %s, want: %s", b, tc.want)
		}
	}
}

func TestJSONPHeader(t *testing.T) {
	var testCases = []struct {
		json     interface{}
		jsonp    bool
		jsonpArg string
		want     string
	}{
		{struct {
			Test int `json:"test"`
		}{1}, true, "testing", headerJSONP},
		{struct {
			Test int `json:"test"`
		}{1}, true, "", headerJSON},
	}

	for _, tc := range testCases {
		var values = url.Values{
			"jsonp": {tc.jsonpArg},
		}
		var URL = &url.URL{
			RawQuery: values.Encode(),
		}
		var rw = &responseWriterTest{h: http.Header{}}
		var c = newHTTPContext(rw, &http.Request{URL: URL})
		c.JSON(tc.json, tc.jsonp)

		if c := rw.h.Get("Content-Type"); c != tc.want {
			t.Errorf("got: %s, want: %s", c, tc.want)
		}
	}
}

func TestSetGet(t *testing.T) {
	var testCases = []struct {
		name  string
		value interface{}
	}{
		{"int", 1},
		{"struct", struct{ x int }{2}},
		{"float", 3.14},
		{"pointer", &(struct{ s int }{1})},
	}

	for _, tc := range testCases {
		var c = newHTTPContext(&responseWriterTest{}, &http.Request{})
		c.Set(tc.name, tc.value)

		if v := c.Get(tc.name); !reflect.DeepEqual(v, tc.value) {
			t.Errorf("got: %s, want: %s", v, tc.value)
		}
	}
}

func TestGetEmpty(t *testing.T) {
	var c = newHTTPContext(&responseWriterTest{}, &http.Request{})

	if v := c.Get("non-existent"); v != nil {
		t.Errorf("got: %s, want: nil", v)
	}
}

func TestParams(t *testing.T) {
	var testCases = []struct {
		name  string
		value string
	}{
		{"id", "123"},
		{"name", "test.txt"},
	}

	for _, tc := range testCases {
		var c = newHTTPContext(&responseWriterTest{}, &http.Request{})
		c.SetParam(tc.name, tc.value)

		if v := c.Param(tc.name); !reflect.DeepEqual(v, tc.value) {
			t.Errorf("got: %s, want: %s", v, tc.value)
		}
	}
}

func TestParamEmpty(t *testing.T) {
	var c = newHTTPContext(&responseWriterTest{}, &http.Request{})

	if v := c.Param("non-existent"); v != "" {
		t.Errorf("got: %s, want: nil", v)
	}
}

func TestRequest(t *testing.T) {
	var r = &http.Request{}
	var c = newHTTPContext(&responseWriterTest{}, r)

	if !reflect.DeepEqual(r, c.Request()) {
		t.Errorf("got: %v, want: %v", c.Request(), r)
	}
}

func TestWrite(t *testing.T) {
	var testCases = []struct {
		content []byte
		code    int
	}{
		{[]byte("Testing"), 200},
	}

	for _, tc := range testCases {
		var rw = &responseWriterTest{}
		var c = newHTTPContext(rw, &http.Request{})

		c.Write(tc.code, tc.content)

		if !reflect.DeepEqual(rw.writtenBytes, tc.content) || tc.code != rw.statusCode {
			t.Errorf("got: %s and %d, want: %s and %d", rw.writtenBytes, rw.statusCode, tc.content, tc.code)
		}
	}
}
