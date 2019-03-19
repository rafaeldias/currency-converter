package currency

// Converter is responsible for querying the convert service.
type Converter interface {
	Convert(from, to string, value float32) (float32, error)
}

// Lister is reponsible for listing the availables currencies to be converted
type Lister interface {
	List() (List, error)
}

// ListConvert groups the basic List and Convert methods of Currency
type ListConverter interface {
	Lister
	Converter
}

type currencyLayer interface {
	Ok() bool
}
