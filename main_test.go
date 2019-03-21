package main

import (
	"os"
	"testing"

	"github.com/rafaeldias/currency-converter/router"
)

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

func TestCurrencyLayerMiddleware(t *testing.T) {
	var hc = &router.HTTPContext{}

	m := currencyLayerMiddleware("x", "y")

	m(hc)

	if c := hc.Get("currency"); c == nil {
		t.Errorf("got: %v, want: currency", c)
	}
}
