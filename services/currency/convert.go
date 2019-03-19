package currency

// Partial payload returned by the Currency Layer, useless properties were omitted
type convertPayload struct {
	Success bool       `json:"success"`
	Result  float32    `json:"result,omitempty"`
	Error   errPayload `json:"error,omitempty"`
}

// convert embeds reference to Credential object in order to execute the requests
type convert struct {
	Credential
}

// Convert requests a conversion from one currency to another with the specified value
func (c *currency) Convert(from, to string, value float32) (float32, error) {
	return 0, nil
}
