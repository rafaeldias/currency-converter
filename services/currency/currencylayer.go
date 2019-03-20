package currency

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type errPayload struct {
	Code int    `json:"code"`
	Info string `json:"info"`
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

type currencyLayer struct {
	host, accessKey string
}

// parseCurrencyLayerPayload reads the stream of bytes from io.Reader and
// decodes it from json to payloader interface, thus being able to
// check if the request was successful without knowing the exact type of
// the payload
func (c *currencyLayer) Parse(p payloader, r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, p); err != nil {
		return err
	}

	return p.Successful()
}

func mountURL(c *currencyLayer, endpoint string, v url.Values) (string, error) {
	// Left the trailing question marking as there will always be at least one
	// query parameter
	var reqURL = bytes.NewBufferString(fmt.Sprintf("%s/api/%s?", c.host, endpoint))

	// Sets access_key, replacing any existing value before it.
	v.Set("access_key", c.accessKey)

	if _, err := reqURL.WriteString(v.Encode()); err != nil {
		return "", err
	}

	return reqURL.String(), nil
}

// Request return an io.ReadCloser or error from the http request to external service.
func (c *currencyLayer) Request(endpoint string, v url.Values) (io.ReadCloser, error) {
	reqURL, err := mountURL(c, endpoint, v)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
