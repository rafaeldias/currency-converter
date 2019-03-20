package currency

// currency implements the ListConverter interface
type currency struct {
	Lister
	Converter
}

// New Returns a ListConverter
func New(host, accessKey string) ListConverter {
	var c = &currencyLayer{host, accessKey}
	return currency{&list{c}, &convert{c}}
}
