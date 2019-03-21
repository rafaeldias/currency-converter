package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rafaeldias/currency-converter/router"
	"github.com/rafaeldias/currency-converter/services/currency"
)

// ValidateConversion validates the user input
func ValidateConversion(hc router.HTTPContexter) {
	if hc.Param("from") == "" || hc.Param("to") == "" {
		hc.Write(http.StatusBadRequest, hc.JSON(Error{"Missing required parameters"}, jsonp))
		hc.Set("valid", false)
		return
	}

	v := hc.Request().URL.Query().Get("value")
	if v == "" {
		hc.Write(http.StatusBadRequest, hc.JSON(Error{"Missing required `value` query parameter"}, jsonp))
		hc.Set("valid", false)
		return
	}

	if _, err := strconv.ParseFloat(v, 32); err != nil {
		hc.Write(http.StatusBadRequest, hc.JSON(Error{fmt.Sprintf("Cannot convert `%s` to float", v)}, jsonp))
		hc.Set("valid", false)
		return

	}

	hc.Set("valid", true)
}

// Conversion handles HTTP requests and convert
// the value using the currency service.
func Conversion(hc router.HTTPContexter) {
	// Invalid data received. Nothing to do from now on
	if valid := hc.Get("valid").(bool); !valid {
		return
	}

	var from = hc.Param("from")
	var to = hc.Param("to")
	var val = hc.Request().URL.Query().Get("value")

	f, _ := strconv.ParseFloat(val, 32)

	log.Printf("Values from query: %s, %s, %f\n", from, to, f)

	// Type assert to currency.Conveter, as this interface
	// has the methods we are interested in for now.
	c := hc.Get("currency").(currency.Converter)

	res, err := c.Convert(from, to, float32(f))
	if err != nil {
		hc.Write(http.StatusInternalServerError, hc.JSON(Error{err.Error()}, jsonp))
		return
	}

	hc.Write(http.StatusOK, hc.JSON(ConversionResult{res}, jsonp))
}
