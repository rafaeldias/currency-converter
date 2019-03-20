package currency

import (
	"fmt"
	"net/url"
)

// convert implements the Converter interface
type convert struct {
	rp requestParser
}

// Convert requests a conversion from one currency to another with the specified value
func (c *convert) Convert(from, to string, value float32) (float32, error) {
	var conv convertPayload
	var v = url.Values{
		"from":   {from},
		"to":     {to},
		"amount": {fmt.Sprintf("%f", value)},
	}

	res, err := c.rp.Request("convert", v)
	if err != nil {
		return 0, err
	}

	defer res.Close()

	if err := c.rp.Parse(&conv, res); err != nil {
		return 0, err
	}

	return conv.Result, nil
}
