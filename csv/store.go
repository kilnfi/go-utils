package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Store allows to store data in csv format
type Store struct {
	path string

	f *os.File

	reader *csv.Reader
	writer *csv.Writer
}

// NewStore creates a Store which data is kept in file path
func NewStore(path string) *Store {
	return &Store{
		path: path,
	}
}

// Open the store file
func (s *Store) Open(createIfNotExist bool) (err error) {
	s.f, err = openFile(s.path, createIfNotExist)
	if err != nil {
		return
	}

	// Set csv reader and writer
	s.reader = csv.NewReader(s.f)
	s.writer = csv.NewWriter(s.f)

	return nil
}

// Reader returns CSV reader
// You MUST call Open() before calling Reader
func (s *Store) Reader() *csv.Reader {
	return s.reader
}

// ReadAll reads all CSV records from disk
func (s *Store) ReadAll() ([][]string, error) {
	err := s.Open(false)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return s.reader.ReadAll()
}

var unmarshallerType = reflect.TypeOf((*Unmarshaller)(nil)).Elem()

// ReadAllStructs reads all CSV records from disk and set it into v
// It assumes v to be a pointer to a slice of object implementing the Unmarshaller interface
func (s *Store) ReadAllStructs(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Slice || rv.Elem().Type().Elem().Kind() != reflect.Ptr || !rv.Elem().Type().Elem().Implements(unmarshallerType) {
		return fmt.Errorf("invalid type %T, expects a slice of %v", v, unmarshallerType)
	}

	records, err := s.ReadAll()
	if err != nil {
		return err
	}

	sv := reflect.MakeSlice(rv.Type().Elem(), 0, len(records))
	for _, rec := range records {
		recV := reflect.New(rv.Type().Elem().Elem().Elem())
		outputs := recV.MethodByName("UnmarshalCSV").Call([]reflect.Value{reflect.ValueOf(rec)})

		if !outputs[0].IsNil() {
			return outputs[0].Interface().(error)
		}
		sv = reflect.Append(sv, recV)
	}

	rv.Elem().Set(sv)

	return nil
}

// Writer returns CSV writer
// You MUST call Open(true) before calling Writer
func (s *Store) Writer() *csv.Writer {
	return s.writer
}

// WriteAll writes all CSV records to disk
func (s *Store) WriteAll(records [][]string) error {
	err := s.Open(true)
	if err != nil {
		return err
	}
	defer s.Close()

	return s.writer.WriteAll(records)
}

var marshallerType = reflect.TypeOf((*Marshaller)(nil)).Elem()

// WriteAllStructs writes all values to disk
// It expects all values to implement Marshaller interface
func (s *Store) WriteAllStructs(values []interface{}) error {
	var records [][]string
	for _, v := range values {
		marshaller, ok := v.(Marshaller)
		if !ok {
			return fmt.Errorf("invalid value type %T does not implement %T", v, marshallerType)
		}
		rec, err := marshaller.MarshalCSV()
		if err != nil {
			return err
		}
		records = append(records, rec)
	}

	return s.WriteAll(records)
}

func (s *Store) Close() error {
	return s.f.Close()
}

func openFile(path string, createIfNotExist bool) (*os.File, error) {
	// check if file exist
	_, err := os.Stat(path)
	if err != nil {
		if !createIfNotExist {
			return nil, fmt.Errorf("file %v does not exist", path)
		}

		// file does not exist so we create it
		err = os.MkdirAll(filepath.Dir(path), 0700)
		if err == nil {
			return os.Create(path)
		}
	} else {
		return os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0644)
	}

	return nil, err
}
