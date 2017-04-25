package fakedata_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

func TestNewColumns(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected fakedata.Columns
	}{
		{name: "one column", input: []string{"email"}, expected: fakedata.Columns{{Key: "email", Name: "email"}}},
		{name: "two columns", input: []string{"email", "domain"}, expected: fakedata.Columns{{Key: "email", Name: "email"}, {Key: "domain", Name: "domain"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := fakedata.NewColumns(tt.input); !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("NewColumns() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestNewColumnsWithSpecs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected fakedata.Columns
	}{
		{name: "one column", input: []string{"login=email"}, expected: fakedata.Columns{{Key: "email", Name: "login"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := fakedata.NewColumns(tt.input); !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("NewColumns() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
