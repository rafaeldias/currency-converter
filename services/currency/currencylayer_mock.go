package currency

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
)

var accessKey = flag.String("accesskey", "abcd", "fake access key for testing purpose")

// Checks for the presence and correctness of access key
func hasAccessKey(r *http.Request) bool {
	return r.URL.Query().Get("access_key") == *accessKey
}

// Creates fake server so we can test if service is making an http request
func currencyLayerServer(l List) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasAccessKey(r) {
			errPay := currenciesPayload{
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
				Success:    true,
				Currencies: l,
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
				Success: true,
				Result:  6.58443,
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
