package currency

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var accessKey = flag.String("accesskey", "abcd", "fake access key for testing purpose")

// Creates fake server so we can test if service is making an http request
func currencyLayerServer(p payloader) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("access_key") != *accessKey {
			res, _ := json.Marshal(p)
			http.Error(w, string(res), http.StatusBadRequest)
			return
		}

		switch r.URL.Path {
		case "/api/list", "/api/convert":
			res, _ := json.Marshal(p)

			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	}))
}

func TestMountURL(t *testing.T) {
	var testCases = []struct {
		host      string
		accessKey string
		endpoint  string
		values    url.Values
		want      string
	}{
		{
			"x",
			"y",
			"z",
			url.Values{
				"q":          {"test"},
				"access_key": {"y"},
			}, "x/api/z?access_key=y&q=test"},
	}

	for _, tc := range testCases {
		c := &currencyLayer{tc.host, tc.accessKey}

		testURL, err := mountURL(c, tc.endpoint, tc.values)
		if err != nil {
			t.Fatalf("Error while testing mountURL: %s", err.Error())
		}

		if testURL != tc.want {
			t.Errorf("got: %s, want: %s", testURL, tc.want)
		}
	}
}

func TestParse(t *testing.T) {
	var testCases = []struct {
		payload, want payloader
	}{
		{
			&currenciesPayload{},
			&currenciesPayload{
				defaultPayload: defaultPayload{
					Success: true,
				},
				Currencies: List{
					"Test": "Testing currency",
				},
			},
		},
		{
			&convertPayload{},
			&convertPayload{
				defaultPayload: defaultPayload{
					Success: true,
				},
				Result: 6.68443,
			},
		},
	}

	for _, tc := range testCases {
		b, err := json.Marshal(tc.want)
		if err != nil {
			t.Fatalf("Error while teting Parse: %s", err.Error())
		}

		r := strings.NewReader(string(b))

		cl := &currencyLayer{}
		if err := cl.Parse(tc.payload, r); err != tc.want.Successful() {
			t.Errorf("got: %s, want: nil", err.Error())
		}

		if !reflect.DeepEqual(tc.want, tc.payload) {
			t.Errorf("got: %v, want: %v", tc.payload, tc.want)
		}
	}
}

func TestParseInvalidJSON(t *testing.T) {
	var r = strings.NewReader("") // Invalid JSON
	var cl = &currencyLayer{}

	if err := cl.Parse(&defaultPayload{}, r); err == nil || err.Error() != "unexpected end of JSON input" {
		t.Errorf("got: %s, want: unexpected end of JSON input", err)
	}
}

type errorReader struct{}

func (er errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("Testing Error")
}

func TestParseReadError(t *testing.T) {
	var r = errorReader{}
	var c = currencyLayer{}

	if err := c.Parse(&defaultPayload{}, r); err == nil || err.Error() != "Testing Error" {
		t.Errorf("got: %s, want: Testing Error", err)
	}
}

func TestRequest(t *testing.T) {
	var testCases = []struct {
		payload, want payloader
		accessKey     string
	}{
		{
			&currenciesPayload{},
			&currenciesPayload{
				defaultPayload: defaultPayload{
					Success: true,
				},
				Currencies: List{"BRA": "Brazil"},
			},
			*accessKey,
		},
		{
			&convertPayload{},
			&convertPayload{
				defaultPayload: defaultPayload{
					Success: true,
				},
				Result: 6.68443,
			},
			*accessKey,
		},
		{
			&defaultPayload{},
			&defaultPayload{
				Success: true,
				Error: errPayload{
					Code: 101,
					Info: "User did not supply an access key or supplied an invalid access key.",
				},
			},
			"",
		},
	}

	for _, tc := range testCases {
		var s = currencyLayerServer(tc.want)

		var c = &currencyLayer{s.URL, tc.accessKey}

		res, err := c.Request("list", url.Values{})
		if err != nil {
			t.Errorf("got: %s, want: nil", err.Error())
		}

		b, err := ioutil.ReadAll(res)
		if err != nil {
			t.Errorf("Error while testing Request method: %s", err.Error())
		}

		if err := json.Unmarshal(b, tc.payload); err != nil {
			t.Errorf("Error while testing Request method: %s", err.Error())
		}

		if !reflect.DeepEqual(tc.want, tc.payload) {
			t.Errorf("got: %s, want: %s", tc.payload, tc.want)
		}

		s.Close()
	}
}

func TestRequestHTTPError(t *testing.T) {
	var s = currencyLayerServer(&defaultPayload{})

	s.Close()

	var c = &currencyLayer{s.URL, ""}

	_, err := c.Request("list", url.Values{})
	if err == nil {
		t.Errorf("got: nil, want: error")
	}
}
