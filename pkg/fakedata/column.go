package fakedata

// A Column represents one field of data to generate
type Column struct {
	Name string
}

// Columns is an array of Column
type Columns []Column

// NewColumns returns an Array of Columns using keys as a specification
func NewColumns(keys []string) (cols Columns) {
	cols = make(Columns, len(keys))

	for i, key := range keys {
		cols[i].Name = key
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
