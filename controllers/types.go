package controllers

const jsonp = true

// Error represents any error message sent by the controller back to the user
type Error struct {
	Message string `json:"message"`
}

// ConversionResult is the payload that will be
// returned to the user by the Conversion function.
type ConversionResult struct {
	Result float32 `json:"result"`
}

// CurrencyList is the payload that will be
// returned to the user by the List function.
type CurrencyList struct {
	Currencies map[string]string `json:"currencies"`
}
