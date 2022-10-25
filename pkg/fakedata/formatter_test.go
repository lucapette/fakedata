package fakedata_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

var columns = fakedata.Columns{{Name: "name", Key: "name"}, {Name: "domain", Key: "domain"}}
var values = []string{"Grace Hopper", "example.com"}

func TestColumnFormatter(t *testing.T) {
	tests := []struct {
		name string
		sep  string
		want string
	}{
		{"default", " ", "Grace Hopper example.com"},
		{"csv", ",", "Grace Hopper,example.com"},
		{"tab", "\t", "Grace Hopper	example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakedata.ColumnFormatter{Separator: tt.sep}
			if got := f.Format(columns, values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ColumnFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLFormatter(t *testing.T) {
	tests := []struct {
		name  string
		table string
		want  string
	}{
		{"table answer", "ANSWER", "INSERT INTO ANSWER (name,domain) VALUES ('Grace Hopper','example.com');"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakedata.SQLFormatter{Table: tt.table}
			if got := f.Format(columns, values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNdjsonFormatter(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"default", "{\"domain\":\"example.com\",\"name\":\"Grace Hopper\"}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakedata.NdjsonFormatter{}
			if got := f.Format(columns, values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NdjsonFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFormatters(b *testing.B) {
	column := &fakedata.ColumnFormatter{}
	b.Run("ColumnFormatter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			column.Format(columns, values)
		}
	})

	ndjson := &fakedata.NdjsonFormatter{}
	b.Run("NdjsonFormatter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ndjson.Format(columns, values)
		}
	})

	sql := &fakedata.SQLFormatter{}
	b.Run("SQLFormatter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sql.Format(columns, values)
		}
	})

}
