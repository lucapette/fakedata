package fakedata_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

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
			args{columns: fakedata.Columns{{Key: "int", Constraints: "10..100"}}, formatter: def},
			1,
			100,
		},
		{
			"int,100..200",
			args{columns: fakedata.Columns{{Key: "int", Constraints: "100..200"}}, formatter: def},
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

func TestGenerateRowWithDateRanges(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		min, max time.Time
	}{
		{
			"date,2016-01-01..2016-12-31",
			args{columns: fakedata.Columns{{Key: "date", Constraints: "2016-01-01..2016-12-31"}}, formatter: def},
			time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2016, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			"date,2016-01-01..",
			args{columns: fakedata.Columns{{Key: "date", Constraints: "2016-01-01"}}, formatter: def},
			time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Now(),
		},
		{
			"date,2046-01-01..2047-01-01",
			args{columns: fakedata.Columns{{Key: "date", Constraints: "2046-01-01..2047-01-01"}}, formatter: def},
			time.Date(2046, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2047, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// this isn't an accurate way of testing random output
			// but it serves a practical purpose
			for index := 0; index < 10000; index++ {
				row := strings.Split(fakedata.GenerateRow(tt.args.columns, tt.args.formatter), " ")

				formattedDate := fmt.Sprintf("%sT00:00:00.000Z", strings.TrimRight(row[0], "\n"))

				actual, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", formattedDate, time.UTC)
				if err != nil {
					t.Fatal(err.Error())
				}

				if !(actual.After(tt.min) && actual.Before(tt.max)) && !actual.Equal(tt.min) && !actual.Equal(tt.max) {
					t.Fatalf("expected a date between %s and %s, but got %s", tt.min, tt.max, actual)
				}
			}
		})
	}
}

func TestGenerateRowWithEnum(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			"enum",
			args{columns: fakedata.Columns{{Key: "enum"}}, formatter: def},
			[]string{"foo", "bar", "baz"},
		},
		{
			"enum,Peter..Olivia..Walter",
			args{columns: fakedata.Columns{{Key: "enum", Constraints: "Peter..Olivia..Walter"}}, formatter: def},
			[]string{"Peter", "Olivia", "Walter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// this isn't an accurate way of testing random output
			// but it serves a practical purpose
			for index := 0; index < 10000; index++ {
				row := strings.TrimRight(fakedata.GenerateRow(tt.args.columns, tt.args.formatter), "\n")

				var found bool
				for _, ex := range tt.expected {
					if strings.Compare(ex, row) == 0 {
						found = true
						break
					}
				}

				if !found {
					t.Fatalf("expected to find %s in %v, but did not", row, tt.expected)
				}
			}
		})
	}
}

func TestValidateGenerators(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"known generator", args{keys: []string{"email", "domain"}}, false},
		{"unknown generator", args{keys: []string{"nogen"}}, true},
		{"mixed generators", args{keys: []string{"nogen", "email", "domain"}}, true},
		{"generator with arguments", args{keys: []string{"int,1..100"}}, false},
		{"mixed generator with arguments", args{keys: []string{"int,1..100", "domain", "email"}}, false},
		{"mixed unknwon generator with arguments", args{keys: []string{"int,1..100", "salery,10k..100k", "email"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fakedata.ValidateGenerators(tt.args.keys); (err != nil) != tt.wantErr {
				t.Errorf("ValidateKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
