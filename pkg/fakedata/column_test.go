package fakedata_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

type columnsTests []struct {
	name     string
	input    []string
	expected fakedata.Columns
	wantErr  bool
}

func TestNewColumns(t *testing.T) {
	tests := columnsTests{
		{
			name:     "one column",
			input:    []string{"email"},
			expected: fakedata.Columns{{Key: "email", Name: "email"}},
			wantErr:  false,
		},
		{
			name:     "two columns",
			input:    []string{"email", "domain"},
			expected: fakedata.Columns{{Key: "email", Name: "email"}, {Key: "domain", Name: "domain"}},
			wantErr:  false,
		},
		{
			name:    "two columns, one column fails",
			input:   []string{"email", "domain", "unsupportedgenerator"},
			wantErr: true,
		},
		{
			name:    "one column, all fails",
			input:   []string{"madeupgenerator"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := fakedata.NewColumns(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("wanted err to be %v but got %v. err: %v", tt.wantErr, (err != nil), err)
			}

			if !tt.wantErr && !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v to equal %v", actual, tt.expected)
			}
		})
	}
}

func TestNewColumnsWithName(t *testing.T) {
	tests := columnsTests{
		{
			name:     "one column",
			input:    []string{"login=email"},
			expected: fakedata.Columns{{Key: "email", Name: "login"}},
			wantErr:  false,
		},
		{
			name:     "one column, unsupported generator",
			input:    []string{"login=notagen"},
			expected: fakedata.Columns{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := fakedata.NewColumns(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("wanted err to be %v but got %v. err: %v", tt.wantErr, (err != nil), err)
			}

			if !tt.wantErr && !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v to equal %v", actual, tt.expected)
			}
		})
	}
}

func TestNewColumnsWithSpec(t *testing.T) {
	tests := columnsTests{
		{
			name:     "int full range",
			input:    []string{"int:1,100"},
			expected: fakedata.Columns{{Key: "int", Name: "int", Options: "1,100"}},
			wantErr:  false,
		},
		{
			name:     "int lower bound",
			input:    []string{"int:1,"},
			expected: fakedata.Columns{{Key: "int", Name: "int", Options: "1,"}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := fakedata.NewColumns(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("wanted err to be %v but got %v. err: %v", tt.wantErr, (err != nil), err)
			}

			if !tt.wantErr && !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v to equal %v", actual, tt.expected)
			}
		})
	}
}
