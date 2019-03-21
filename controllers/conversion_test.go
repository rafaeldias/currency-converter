package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/rafaeldias/currency-converter/services/currency"
)

type converterTest struct {
	err    error
	called bool
}

func (c *converterTest) Convert(from, to string, value float32) (float32, error) {
	c.called = true
	return 0, c.err
}

type ctxHTTPConversionTest struct {
	converter    currency.Converter
	request      *http.Request
	params       map[string]string
	data         map[string]interface{}
	code         int
	bytes        []byte
	writtenBytes []byte
	output       interface{}
}

func (c *ctxHTTPConversionTest) JSON(v interface{}, jsonp bool) []byte {
	c.output = v

	return c.writtenBytes
}

func (c *ctxHTTPConversionTest) Get(s string) interface{} {
	return c.data[s]
}

func (c *ctxHTTPConversionTest) Set(s string, v interface{}) {
	c.data[s] = v
}

func (c *ctxHTTPConversionTest) Request() *http.Request {
	return c.request
}

func (c *ctxHTTPConversionTest) SetParam(name, value string) {}

func (c *ctxHTTPConversionTest) Param(name string) string {
	return c.params[name]
}

func (c *ctxHTTPConversionTest) Write(code int, b []byte) error {
	c.bytes = b
	c.code = code
	return nil
}

func (c *ctxHTTPConversionTest) Header() http.Header {
	return nil
}

var testCases = []struct {
	params map[string]string
	req    *http.Request
	want   bool
}{
	{
		map[string]string{"from": "x", "to": "x"},
		&http.Request{
			URL: &url.URL{
				RawQuery: "value=10",
			},
		},
		true,
	},
	{
		map[string]string{"from": "", "to": "x"},
		&http.Request{
			URL: &url.URL{
				RawQuery: "value=10",
			},
		},
		false,
	},
	{
		map[string]string{"from": "x", "to": ""},
		&http.Request{
			URL: &url.URL{
				RawQuery: "value=10",
			},
		},
		false,
	},
	{
		map[string]string{"from": "x", "to": "x"},
		&http.Request{
			URL: &url.URL{},
		},
		false,
	},
	{
		map[string]string{"from": "x", "to": ""},
		&http.Request{
			URL: &url.URL{
				RawQuery: "value=a",
			},
		},
		false,
	},
}

func TestValidateConversionValid(t *testing.T) {
	for _, tc := range testCases {
		ctx := &ctxHTTPConversionTest{
			params:  tc.params,
			request: tc.req,
			data:    map[string]interface{}{},
		}

		ValidateConversion(ctx)

		if v := ctx.Get("valid"); !reflect.DeepEqual(v, tc.want) {
			t.Errorf("got: %v, want: %v", v, tc.want)
		}
	}
}

func TestValidateConversionJSON(t *testing.T) {
	for _, tc := range testCases[1:] {
		ctx := &ctxHTTPConversionTest{
			params:  tc.params,
			request: tc.req,
			data:    map[string]interface{}{},
		}

		ValidateConversion(ctx)

		if ctx.output == nil {
			t.Error("got: nil, want: error")
		}
	}
}

func TestValidateConversionStatus(t *testing.T) {
	for _, tc := range testCases[1:] {
		ctx := &ctxHTTPConversionTest{
			params:  tc.params,
			request: tc.req,
			data:    map[string]interface{}{},
		}

		ValidateConversion(ctx)

		if ctx.code != http.StatusBadRequest {
			t.Errorf("got: %d, want: %d", ctx.code, http.StatusBadRequest)
		}
	}
}

func TestValidateConversionWrite(t *testing.T) {
	for _, tc := range testCases[1:] {
		var want = []byte("x")
		ctx := &ctxHTTPConversionTest{
			params:       tc.params,
			request:      tc.req,
			data:         map[string]interface{}{},
			writtenBytes: want,
		}

		ValidateConversion(ctx)

		if !reflect.DeepEqual(ctx.bytes, want) {
			t.Errorf("got: %s, want: %s", ctx.bytes, want)
		}
	}
}

func TestConversionValid(t *testing.T) {
	var c = &converterTest{}
	var testCases = []struct {
		params map[string]string
		req    *http.Request
		data   map[string]interface{}
		want   bool
	}{
		{
			map[string]string{"from": "x", "to": "x"},
			&http.Request{
				URL: &url.URL{
					RawQuery: "value=10",
				},
			},
			map[string]interface{}{
				"valid":    false,
				"currency": c,
			},
			false,
		},
		{
			map[string]string{"from": "x", "to": "x"},
			&http.Request{
				URL: &url.URL{
					RawQuery: "value=10",
				},
			},
			map[string]interface{}{
				"valid":    true,
				"currency": c,
			},
			true,
		},
	}

	for _, tc := range testCases {
		ctx := &ctxHTTPConversionTest{
			converter: c,
			params:    tc.params,
			request:   tc.req,
			data:      tc.data,
		}

		Conversion(ctx)

		if c.called != tc.want {
			t.Errorf("got: %v, want: %v", c.called, tc.want)
		}
	}
}

func TestConversionStatus(t *testing.T) {
	var testCases = []struct {
		params map[string]string
		req    *http.Request
		data   map[string]interface{}
		want   int
	}{
		{
			map[string]string{"from": "x", "to": "x"},
			&http.Request{
				URL: &url.URL{
					RawQuery: "value=10",
				},
			},
			map[string]interface{}{
				"valid":    true,
				"currency": &converterTest{},
			},
			http.StatusOK,
		},
		{
			map[string]string{"from": "x", "to": "x"},
			&http.Request{
				URL: &url.URL{
					RawQuery: "value=10",
				},
			},
			map[string]interface{}{
				"valid":    true,
				"currency": &converterTest{err: errors.New("")},
			},
			http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		ctx := &ctxHTTPConversionTest{
			params:  tc.params,
			request: tc.req,
			data:    tc.data,
		}

		Conversion(ctx)

		if ctx.code != tc.want {
			t.Errorf("got: %v, want: %v", ctx.code, tc.want)
		}
	}
}
