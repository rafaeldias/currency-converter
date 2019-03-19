package currency

// currency implements the ListConverter interface
type currency struct {
	Lister
	Converter
}

// New Returns a ListConverter
func New(c Credential) ListConverter {
	return &currency{&list{c}, &convert{c}}
}
