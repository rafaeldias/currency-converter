package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// convert embeds reference to Credential object in order to execute the requests
type convert struct {
	Credential
}

// Convert requests a conversion from one currency to another with the specified value
func (c *convert) Convert(from, to string, value float32) (float32, error) {
	var conv convertPayload
	var convURL = fmt.Sprintf("%s/api/convert?from=%s&to=%s&amount=%f&access_key=%s", c.Host, from, to, value, c.AccessKey)

	res, err := http.Get(convURL)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	converted, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(converted, &conv); err != nil {
		return 0, err
	}

	if !conv.Success {
		return 0, errors.New(conv.Error.Info)
	}

	return conv.Result, nil
}
