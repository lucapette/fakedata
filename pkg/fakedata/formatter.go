package fakedata

import (
	"bytes"
	"fmt"
	"strings"
)

// Formatter is the interface wraps the Format method we use to format each row
type Formatter interface {
	Format(Columns, []string) string
}

// SeparatorFormatter is a Formatter for character separated formats
type SeparatorFormatter struct {
	Separator string
}

// Format as character separated strings
func (f *SeparatorFormatter) Format(columns Columns, values []string) string {
	return strings.Join(values, f.Separator)
}

// SQLFormatter is a Formatter for the SQL insert statement
type SQLFormatter struct {
	Table string
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

	fmt.Fprintf(sql, // nolint: errcheck
		"INSERT INTO %s (%s) VALUES (%s);",
		f.Table,
		strings.Join(names, ","),
		strings.Join(formattedValues, ","),
	)

	return sql.String()
}

// NewSeparatorFormatter returns a SeparatorFormatter using the sep string as a separator
func NewSeparatorFormatter(sep string) (f *SeparatorFormatter) {
	return &SeparatorFormatter{Separator: sep}
}

// NewSQLFormatter returns a SQLFormatter using the table string for table name generation
func NewSQLFormatter(table string) (f *SQLFormatter) {
	return &SQLFormatter{Table: table}
}
