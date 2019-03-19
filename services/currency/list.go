package currency

import (
	_ "encoding/json"
	_ "errors"
	"fmt"
	_ "io/ioutil"
	"net/http"
)

// List is the collection of available currencies to be converted
type List map[string]string

// list embeds reference to Credential object in order to execute the requests
type list struct {
	Credential
}

// List requests a collection of currencies from the external service
func (l *list) List() (List, error) {
	var list currenciesPayload
	var listURL = fmt.Sprintf("%s/api/list?access_key=%s", l.Host, l.AccessKey)

	res, err := http.Get(listURL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if err := parseCurrencyLayerPayload(&list, res.Body); err != nil {
		return nil, err
	}

	return list.Currencies, nil
}
