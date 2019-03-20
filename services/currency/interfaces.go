package currency

import (
	"io"
	"net/url"
)

// Converter is responsible for querying the convert service.
type Converter interface {
	Convert(from, to string, value float32) (float32, error)
}

// Lister is reponsible for listing the availables currencies to be converted
type Lister interface {
	List() (List, error)
}

// ListConverter groups the basic List and Convert methods of Currency
type ListConverter interface {
	Lister
	Converter
}

// payloader is a private interface  used to get the status of the request
type payloader interface {
	Successful() error
}

// requestPayload requests the endpoint with the url.Values
type requester interface {
	Request(endpoint string, v url.Values) (io.ReadCloser, error)
}

// parser decodes io.Reader bytes to payloader interface
type parser interface {
	Parse(p payloader, r io.Reader) error
}

// requestParser groups the external request parser behaviour
type requestParser interface {
	requester
	parser
}
