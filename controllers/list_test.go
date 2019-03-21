package controllers

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/rafaeldias/currency-converter/services/currency"
)

type listerTest struct {
	err error
}

func (l *listerTest) List() (currency.List, error) {
	return currency.List{}, l.err
}

type ctxHTTPListTest struct {
	lister       currency.Lister
	code         int
	bytes        []byte
	writtenBytes []byte
	output       interface{}
}

func (c *ctxHTTPListTest) JSON(v interface{}, jsonp bool) []byte {
	c.output = v

	return c.writtenBytes
}

func (c *ctxHTTPListTest) Get(s string) interface{} {
	return c.lister
}

func (c *ctxHTTPListTest) Set(s string, v interface{}) {}

func (c *ctxHTTPListTest) Request() *http.Request {
	return nil
}

func (c *ctxHTTPListTest) SetParam(name, value string) {}

func (c *ctxHTTPListTest) Param(name string) string {
	return ""
}

func (c *ctxHTTPListTest) Write(code int, b []byte) error {
	c.bytes = b
	c.code = code
	return nil
}

func (c *ctxHTTPListTest) Header() http.Header {
	return nil
}

func TestListJSON(t *testing.T) {
	var testCases = []struct {
		err  error
		want interface{}
	}{
		{nil, &CurrencyList{currency.List{}}},
		{errors.New("List Write Error"), &Error{"List Write Error"}},
	}

	for _, tc := range testCases {
		ctx := &ctxHTTPListTest{lister: &listerTest{tc.err}}

		List(ctx)

		if reflect.DeepEqual(ctx.output, tc.want) {
			t.Errorf("got: %v, want: %v", ctx.output, tc.want)
		}
	}
}

func TestListStatus(t *testing.T) {
	var testCases = []struct {
		err  error
		want int
	}{
		{nil, http.StatusOK},
		{errors.New("List Write Error"), http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		ctx := &ctxHTTPListTest{lister: &listerTest{tc.err}}

		List(ctx)

		if ctx.code != tc.want {
			t.Errorf("got: %v, want: %v", ctx.code, tc.want)
		}
	}
}

func TestListWriteBytes(t *testing.T) {
	var testCases = []struct {
		err  error
		want []byte
	}{
		{nil, []byte("x")},
		{errors.New("List Write Error"), []byte("y")},
	}

	for _, tc := range testCases {
		ctx := &ctxHTTPListTest{writtenBytes: tc.want, lister: &listerTest{tc.err}}

		List(ctx)

		if !reflect.DeepEqual(ctx.bytes, tc.want) {
			t.Errorf("got: %s, want: %s", ctx.bytes, tc.want)
		}
	}
}
