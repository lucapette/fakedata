package fakedata

import (
	"bytes"
	"fmt"
	"strings"
)

var unknownGeneratorError = `
  Unknown generator: %s.`

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

// ValidateGenerators validates the passed generators agains available generators
func ValidateGenerators(keys []string) error {
	// check each parameter
	var validationError []string
	var err error
	for _, k := range keys {
		paramArg := strings.Split(k, ",")
		// Seperate arguments that have parameters, e.g. int,1..50
		if len(paramArg) > 1 {
			k = paramArg[0]
		}
		if generators[k].Name != k {
			validationError = append(validationError, fmt.Errorf(unknownGeneratorError, k).Error())
		}
	}
	// If there are errors, join them into one big string
	if len(validationError) > 0 {
		err = fmt.Errorf(strings.Join(validationError, ""))
	}
	return err
}
