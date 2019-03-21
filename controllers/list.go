package controllers

import (
	"net/http"

	"github.com/rafaeldias/currency-converter/router"
	"github.com/rafaeldias/currency-converter/services/currency"
)

// List handles HTTP requests and convert
// the value using the currency service.
func List(hc router.HTTPContexter) {
	// Type assert to currency.Lister, as this interface
	// has the methods we are interested in for now.
	l := hc.Get("currency").(currency.Lister)

	res, err := l.List()
	if err != nil {
		hc.Write(http.StatusInternalServerError, hc.JSON(Error{err.Error()}, jsonp))
		return
	}

	hc.Write(http.StatusOK, hc.JSON(CurrencyList{res}, jsonp))
}
