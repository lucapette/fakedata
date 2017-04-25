package fakedata

import (
	"bytes"
	"sort"
)

// GenerateRow generates a row of fake data using Columns
// in the specified format
func GenerateRow(columns Columns, formatter Formatter) string {
	var output bytes.Buffer

	genValues := make([]string, len(columns))
	for i, column := range columns {
		genValues[i] = generate(column)
	}

	output.WriteString(formatter.Format(columns, genValues))

	output.WriteString("\n")

	return output.String()
}

// Generators returns all the available generators
func Generators() []string {
	gens := make([]string, 0)

	for k := range generators {
		gens = append(gens, k)
	}

	sort.Strings(gens)
	return gens
}
