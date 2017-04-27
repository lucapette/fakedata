package fakedata_test

import (
	"strconv"
	"strings"
	"testing"

	"regexp"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

var csv = fakedata.NewSeparatorFormatter(",")
var def = fakedata.NewSeparatorFormatter(" ")
var tab = fakedata.NewSeparatorFormatter("\t")

type args struct {
	columns   fakedata.Columns
	formatter fakedata.Formatter
}

func TestGenerateRow(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{"email", args{columns: fakedata.Columns{{Key: "email"}}, formatter: def}, `^.+?@.+?\..+$`},
		{"domain", args{columns: fakedata.Columns{{Key: "domain"}}, formatter: def}, `^.+?\..+?$`},
		{"username", args{columns: fakedata.Columns{{Key: "username"}}, formatter: def}, `^[a-zA-Z0-9]{2,}$`},
		{"double", args{columns: fakedata.Columns{{Key: "double"}}, formatter: def}, `^-?[0-9]+?(\.[0-9]+?)?$`},
		{"date", args{columns: fakedata.Columns{{Key: "date"}}, formatter: def}, `^\d{4}-\d{2}-\d{2}$`},
		{"username domain", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, formatter: def}, `^[a-zA-Z0-9]{2,} .+?\..+?$`},
		{"username domain csv", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, formatter: csv}, `^[a-zA-Z0-9]{2,},.+?\..+?$`},
		{"username domain tab", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, formatter: tab}, `^[a-zA-Z0-9]{2,}\t.+?\..+?$`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := strings.TrimRight(fakedata.GenerateRow(tt.args.columns, tt.args.formatter), "\n")

			matched, err := regexp.MatchString(tt.expected, actual)
			if err != nil {
				t.Error(err.Error())
			}

			if !matched {
				t.Errorf("expected %s to match '%s', but did not", tt.expected, actual)
			}
		})
	}
}

func TestGenerateRowWithIntRanges(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		min, max int
	}{
		{
			"int,1..10",
			args{columns: fakedata.Columns{{Key: "int", Min: "10", Max: "100"}}, formatter: def},
			1,
			100,
		},
		{
			"int,100..200",
			args{columns: fakedata.Columns{{Key: "int", Min: "100", Max: "200"}}, formatter: def},
			100,
			1200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// this isn't an accurate way of testing random output
			// but it serves a practical purpose
			for index := 0; index < 10000; index++ {
				row := strings.Split(fakedata.GenerateRow(tt.args.columns, tt.args.formatter), " ")
				actual, err := strconv.Atoi(strings.TrimRight(row[0], "\n"))
				if err != nil {
					t.Fatal(err.Error())
				}

				if !(actual >= tt.min && actual <= tt.max) {
					t.Fatalf("expected a number between %d and %d, but got %d", tt.min, tt.max, actual)
				}
			}
		})
	}
}
