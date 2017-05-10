package fakedata

import (
	"fmt"
	"strings"
)

// A Column represents one field of data to generate
type Column struct {
	Name        string
	Key         string
	Constraints string
}

func (c *Column) String() string {
	return fmt.Sprintf("Column(%s=%s)", c.Name, c.Key)
}

// Columns is an array of Column
type Columns []Column

// NewColumns returns an array of Columns using keys as a specification
func NewColumns(keys []string) (cols Columns) {
	cols = make(Columns, len(keys))

	for i, k := range keys {
		specs := strings.Split(k, ",")

		if len(specs) > 1 {
			cols[i].Constraints = specs[1]
		}

		values := strings.Split(specs[0], "=")
		var name, key string

		if len(values) == 2 {
			name = values[0]
			key = values[1]
		} else {
			name = values[0]
			key = values[0]
		}

		cols[i].Name = name
		cols[i].Key = key
	}

	return cols
}

func (columns Columns) names() (names []string) {
	names = make([]string, len(columns))

	for i, field := range columns {
		names[i] = field.Name
	}

	return names
}
