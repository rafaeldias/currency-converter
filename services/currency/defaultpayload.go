package currency

import "errors"

// Commom properties of Currency Layer payload
type defaultPayload struct {
	Success bool       `json:"success"`
	Error   errPayload `json:"error,omitempty"`
}

// Suceeded checks if payload was successfully received
func (d *defaultPayload) Successful() error {
	if !d.Success {
		return errors.New(d.Error.Info)
	}
	return nil
}
