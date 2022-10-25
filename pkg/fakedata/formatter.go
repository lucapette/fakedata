package fakedata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Formatter is the interface wraps the Format method we use to format each row
type Formatter interface {
	Format(Columns, []string) string
}

// ColumnFormatter is a Formatter for character separated formats
type ColumnFormatter struct {
	Separator string
}

// SQLFormatter is a Formatter for the SQL insert statement
type SQLFormatter struct {
	Table string
}

// NdJsonFormatter is a Formatter for http://ndjson.org/
type NdjsonFormatter struct {
}

// Format as character separated strings
func (f *ColumnFormatter) Format(columns Columns, values []string) string {
	return strings.Join(values, f.Separator)
}

// Format as SQL statements
func (f *SQLFormatter) Format(columns Columns, values []string) string {
	sql := &bytes.Buffer{}
	names := make([]string, len(columns))

	for i, field := range columns {
		names[i] = field.Name
	}

	formattedValues := make([]string, len(columns))
	for i, value := range values {
		formattedValues[i] = fmt.Sprintf("'%s'", value)
	}

	fmt.Fprintf(sql,
		"INSERT INTO %s (%s) VALUES (%s);",
		f.Table,
		strings.Join(names, ","),
		strings.Join(formattedValues, ","),
	)

	return sql.String()
}

// Format as ndjson
func (f *NdjsonFormatter) Format(columns Columns, values []string) string {
	data := make(map[string]string, len(columns))

	for i := 0; i < len(columns); i++ {
		data[columns[i].Name] = values[i]
	}

	v, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(v)
}

// NewColumnFormatter returns a ColumnFormatter using the sep string as a separator
func NewColumnFormatter(sep string) (f *ColumnFormatter) {
	return &ColumnFormatter{Separator: sep}
}

// NewSQLFormatter returns a SQLFormatter using the table string for table name generation
func NewSQLFormatter(table string) (f *SQLFormatter) {
	return &SQLFormatter{Table: table}
}

// NewNdjsonFormatter returns a NdjsonFormatter
func NewNdjsonFormatter() (f *NdjsonFormatter) {
	return &NdjsonFormatter{}
}
