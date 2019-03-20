package currency

import "net/url"

// List is the collection of available currencies to be converted
type List map[string]string

// list implements the Lister interface
type list struct {
	rp requestParser
}

// List requests a collection of currencies from the external service
func (l *list) List() (List, error) {
	var list currenciesPayload

	res, err := l.rp.Request("list", url.Values{})
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if err := l.rp.Parse(&list, res); err != nil {
		return nil, err
	}

	return list.Currencies, nil
}
