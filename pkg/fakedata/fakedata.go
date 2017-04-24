package fakedata

import (
	"bytes"
	"strings"
)

func joinFunc(sep string) func([]string) string {
	return func(keys []string) string {
		return strings.Join(keys, sep)
	}
}

func formatter(format string) (f func([]string) string) {
	switch format {
	case "tab":
		f = joinFunc("\t")
	case "csv":
		f = joinFunc(",")
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

// GenerateRow generates a row of fake data according to key
// in the specified format
func GenerateRow(keys []string, format string) string {
	var output bytes.Buffer

	f := formatter(format)

	values := make([]string, len(keys))

	for i, k := range keys {
		values[i] = generate(k)
	}

	output.WriteString(f(values))

	output.WriteString("\n")

	return output.String()
}

// List returns all the available generators
func List() []string {
	list := make([]string, 0)

	for k := range generators {
		list = append(list, k)
	}
	return list
}
