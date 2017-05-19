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
		wantErr  bool
	}{
		{name: "one column", input: []string{"email"}, expected: fakedata.Columns{{Key: "email", Name: "email"}}, wantErr: false},
		{name: "two columns", input: []string{"email", "domain"}, expected: fakedata.Columns{{Key: "email", Name: "email"}, {Key: "domain", Name: "domain"}}, wantErr: false},
		{name: "two columns, one column fails", input: []string{"email", "domain", "unsupportedgenerator"}, expected: nil, wantErr: true},
		{name: "one column, all fails", input: []string{"madeupgenerator"}, expected: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual, err := fakedata.NewColumns(tt.input); !reflect.DeepEqual(actual, tt.expected) && ((err != nil) != tt.wantErr) {
				t.Errorf("NewColumns() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestNewColumnsWithName(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected fakedata.Columns
		wantErr  bool
	}{
		{name: "one column", input: []string{"login=email"}, expected: fakedata.Columns{{Key: "email", Name: "login"}}, wantErr: false},
		{name: "one column, unupported generator", input: []string{"login=notagen"}, expected: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual, err := fakedata.NewColumns(tt.input); !reflect.DeepEqual(actual, tt.expected) && ((err != nil) != tt.wantErr) {
				t.Errorf("NewColumns() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestNewColumnsWithSpec(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected fakedata.Columns
		wantErr  bool
	}{
		{name: "int full range", input: []string{"int,1..100"}, expected: fakedata.Columns{{Key: "int", Name: "int", Constraints: "1..100"}}, wantErr: false},
		{name: "int lower bound", input: []string{"int,1.."}, expected: fakedata.Columns{{Key: "int", Name: "int", Constraints: "1.."}}, wantErr: false},
		{name: "int lower bound no range syntax", input: []string{"int,10"}, expected: fakedata.Columns{{Key: "int", Name: "int", Constraints: "10"}}, wantErr: false},
		{name: "int spelled wrong", input: []string{"integer,10"}, expected: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual, err := fakedata.NewColumns(tt.input); !reflect.DeepEqual(actual, tt.expected) && ((err != nil) != tt.wantErr) {
				t.Errorf("NewColumns() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
