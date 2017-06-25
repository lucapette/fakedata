package fakedata

import (
	"bytes"
	"fmt"
	"strings"
)

// A Column represents one field of data to generate
type Column struct {
	Name     string
	Key      string
	Generate func() string
}

// Columns is an array of Column
type Columns []Column

// NewColumns returns an array of Columns using keys as a specification.
// It returns an error with a line for each unknown key
func NewColumns(keys []string) (cols Columns, err error) {
	cols = make(Columns, len(keys))

	f := newFactory()

	for i, k := range keys {
		specs := strings.Split(k, ":")

		values := strings.Split(specs[0], "=")
		var name, key, options string

		if len(values) == 2 {
			name = values[0]
			key = values[1]
		} else {
			name = values[0]
			key = values[0]
		}

		if len(specs) > 1 {
			options = specs[1]
		}

		fn, err := f.extractFunc(key, options)
		if err != nil {
			return cols, err
		}

		cols[i].Name = name
		cols[i].Key = key
		cols[i].Generate = fn
	}

	return cols, err
}

// GenerateRow generates a row of fake data using columns
// in the specified format
func (columns Columns) GenerateRow(formatter Formatter) string {
	row := &bytes.Buffer{}

	values := make([]string, len(columns))
	for i, column := range columns {
		values[i] = column.Generate()
	}

	fmt.Fprintf(row, "%s", formatter.Format(columns, values))

	return row.String()
}
