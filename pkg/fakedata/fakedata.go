package fakedata

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func joinFunc(sep string) func(Columns, []string) string {
	return func(columns Columns, values []string) string {
		return strings.Join(values, sep)
	}
}

func sqlFunc() func(Columns, []string) string {
	return func(columns Columns, values []string) string {
		sql := bytes.NewBufferString("INSERT INTO TABLE (")

		sql.WriteString(strings.Join(columns.names(), ","))
		sql.WriteString(") values (")

		formattedValues := make([]string, len(columns))
		for i, value := range values {
			formattedValues[i] = fmt.Sprintf("'%s'", value)
		}
		sql.WriteString(strings.Join(formattedValues, ","))

		sql.WriteString(");")
		return sql.String()
	}
}

func formatter(format string) (f func(Columns, []string) string) {
	switch format {
	case "tab":
		f = joinFunc("\t")
	case "csv":
		f = joinFunc(",")
	case "sql":
		f = sqlFunc()
	default:
		f = joinFunc(" ")
	}
	return f
}

func generate(key string) string {
	if f, ok := generators[key]; ok {
		return f()
	}

	return ""
}

// GenerateRow generates a row of fake data using Columns
// in the specified format
func GenerateRow(columns Columns, format string) string {
	var output bytes.Buffer

	f := formatter(format)

	genValues := make([]string, len(columns))
	for i, field := range columns {
		genValues[i] = generate(field.Name)
	}

	output.WriteString(f(columns, genValues))

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
