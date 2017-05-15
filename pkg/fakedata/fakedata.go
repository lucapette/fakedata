package fakedata

import (
	"bytes"
	"fmt"
	"strings"
)

// GenerateRow generates a row of fake data using columns
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

// ValidateGenerators validates each key in keys against available generators
func ValidateGenerators(keys []string) (err error) {
	var errors bytes.Buffer

	for _, k := range keys {
		key := strings.Split(k, ",")[0]

		if _, ok := generators[key]; !ok {
			errors.WriteString(fmt.Sprintf("Unknown generator: %s.\n", key))
		}
	}

	if errors.Len() > 0 {
		err = fmt.Errorf(errors.String())
	}
	return err
}
