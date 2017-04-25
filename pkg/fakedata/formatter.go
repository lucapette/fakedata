package fakedata

import (
	"bytes"
	"fmt"
	"strings"
)

// Formatter is that wraps the Format method we use to format each row
type Formatter interface {
	Format(Columns, []string) string
}

// SeparatorFormatter is a Formatter for characther separated formats
type SeparatorFormatter struct {
	Separator string
}

// Format as characther separated strings
func (f *SeparatorFormatter) Format(columns Columns, values []string) string {
	return strings.Join(values, f.Separator)
}

// SQLFormatter is a Formatter for the SQL insert statement
type SQLFormatter struct {
	Table string
}

// Format as SQL statements
func (f *SQLFormatter) Format(columns Columns, values []string) string {
	sql := bytes.NewBufferString(fmt.Sprintf("INSERT INTO %s (", f.Table))

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

// NewSeparatorFormatter returns a SeparatorFormatter using the sep string as a separator
func NewSeparatorFormatter(sep string) (f *SeparatorFormatter) {
	return &SeparatorFormatter{Separator: sep}

}

// NewSQLFormatter returns a SQLFormatter using the table string for table name generation
func NewSQLFormatter(table string) (f *SQLFormatter) {
	return &SQLFormatter{Table: table}
}
