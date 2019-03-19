package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// List is the collection of available currencies to be converted
type List map[string]string

// Partial payload returned by the Currency Layer, useless properties were omitted
type currenciesPayload struct {
	Success    bool       `json:"success"`
	Currencies List       `json:"currencies,omitempty"`
	Error      errPayload `json:"error,omitempty"`
}

// list embeds reference to Credential object in order to execute the requests
type list struct {
	Credential
}

// List requests a collection of currencies from the external service
func (l *list) List() (List, error) {
	var list currenciesPayload

	res, err := http.Get(fmt.Sprintf("%s/api/list?access_key=%s", l.Host, l.AccessKey))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	currencies, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(currencies, &list); err != nil {
		return nil, err
	}

	if !list.Success {
		return nil, errors.New(list.Error.Info)
	}

	return list.Currencies, nil
}
