package currency

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var accessKey = flag.String("accesskey", "abcd", "fake access key for testing purpose")

type errPayload struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// Commom properties of Currency Layer payloaa
type defaultPayload struct {
	Success bool       `json:"success"`
	Error   errPayload `json:"error,omitempty"`
}

func (d *defaultPayload) Ok() bool {
	return d.Success
}

// Partial payload returned by the Currency Layer, useless properties were omitted
type currenciesPayload struct {
	defaultPayload

	Currencies List `json:"currencies,omitempty"`
}

// Partial payload returned by the Currency Layer, useless properties were omitted
type convertPayload struct {
	defaultPayload

	Result float32 `json:"result,omitempty"`
}

func parseCurrencyLayerPayload(c currencyLayer, r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, c); err != nil {
		return err
	}

	if !c.Ok() {
		return errors.New(list.Error.Info)
	}

	return nil
}

// Checks for the presence and correctness of access key
func hasAccessKey(r *http.Request) bool {
	return r.URL.Query().Get("access_key") == *accessKey
}

// Creates fake server so we can test if service is making an http request
func currencyLayerServer(l List) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasAccessKey(r) {
			errPay := defaultPayload{
				Success: false,
				Error: errPayload{
					Code: http.StatusSwitchingProtocols,
					Info: "User did not supply an access key or supplied an invalid access key.",
				},
			}

			res, err := json.Marshal(errPay)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(res), http.StatusBadRequest)
			return
		}

		switch r.URL.Path {
		case "/api/list":
			currPay := currenciesPayload{
				defaultPayload{
					Success: true,
				},
				l,
			}

			res, err := json.Marshal(currPay)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(res)
		case "/api/convert":
			convPay := convertPayload{
				defaultPayload{
					Success: true,
				},
				6.58443,
			}

			res, err := json.Marshal(convPay)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(res)

		default:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}))
}
