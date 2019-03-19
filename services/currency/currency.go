package currency

// errPayload represents the entity returned by currency Layer when an error occurs
type errPayload struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// currency embeds list and convert in order to match the ListConvert Interface
type currency struct {
	*list
	*convert
}

// New receives the Credentials and returns an instance of ListConvert interface
func New(c Credential) ListConvert {
	return &currency{&list{c}, &convert{c}}
}
