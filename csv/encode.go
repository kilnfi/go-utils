package csv

// Marshaller is the interface implemented by types that can marshal themselves into valid CSV record.
type Marshaller interface {
	MarshalCSV() ([]string, error)
}

// Unmarshaler is the interface implemented by types that can unmarshal a CSV record description of themselves.
type Unmarshaller interface {
	UnmarshalCSV([]string) error
}
