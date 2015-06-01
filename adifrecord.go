package adifparser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Public interface for ADIFRecords
type ADIFRecord interface {
	// Print as ADIF String
	ToString() string
	// Fingerprint for duplication detection
	Fingerprint() string
	// Setters and getters
}

// Internal implementation for ADIFRecord
type baseADIFRecord struct {
	values map[string]string
}

type fieldData struct {
	name     string
	value    string
	typecode byte
	hasType  bool
}

// Create a new ADIFRecord from scratch
func NewADIFRecord() *baseADIFRecord {
	record := &baseADIFRecord{}
	record.values = make(map[string]string)
	return record
}

// Parse an ADIFRecord
func ParseADIFRecord(buf []byte) (*baseADIFRecord, error) {
	record := NewADIFRecord()

	for len(buf) > 0 {
		var data *fieldData
		var err error
		data, buf, err = getNextField(buf)
		if err != nil {
			return nil, err
		}
		// TODO: accomodate types
		record.values[data.name] = data.value
	}

	return record, nil
}

// Get the next field, return field data, leftover data, and optional error
func getNextField(buf []byte) (*fieldData, []byte, error) {
	data := &fieldData{}

	// Extract name
	start_of_name := bytes.IndexByte(buf, '<') + 1
	end_of_name := bytes.IndexByte(buf, ':')
	data.name = strings.ToLower(string(buf[start_of_name:end_of_name]))
	buf = buf[end_of_name+1:]

	// Length
	var length int
	var err error
	end_of_tag := bytes.IndexByte(buf, '>')
	start_type := bytes.IndexByte(buf, ':')
	if start_type == -1 || start_type > end_of_tag {
		end_of_length := bytes.IndexByte(buf, '>')
		length, err = strconv.Atoi(string(buf[:end_of_length]))
		buf = buf[end_of_length+1:]
		data.hasType = false
	} else {
		length, err = strconv.Atoi(string(buf[:start_type]))
		data.typecode = buf[start_type+1]
		data.hasType = true
		buf = buf[start_type+3:]
	}
	if err != nil {
		// TODO: log the error
		return nil, buf, err
	}

	// Value
	data.value = string(buf[:length])
	buf = bytes.TrimSpace(buf[length:])

	return data, buf, nil
}

func serializeField(name string, value string) string {
	return fmt.Sprintf("<%s:%d>%s", name, len(value), value)
}

// Print an ADIFRecord as a string
func (r *baseADIFRecord) ToString() string {
	var record bytes.Buffer
	for n, v := range r.values {
		record.WriteString(serializeField(n, v))
	}
	return record.String()
}

// Get fingerprint of ADIFRecord
func (r *baseADIFRecord) Fingerprint() string {
	return ""
}

// Get a value
func (r *baseADIFRecord) GetValue(name string) string {
	return r.values[name]
}
