package fakedata_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

type columnsTests []struct {
	name     string
	input    []string
	expected fakedata.Columns
	wantErr  bool
}

func columnsEql(a, b fakedata.Columns) bool {
	if len(a) != len(b) {
		return false
	}

	for i, col := range a {
		if !(reflect.DeepEqual(col.Name, b[i].Name) &&
			reflect.DeepEqual(col.Key, b[i].Key)) {
			return false
		}
	}
	return true
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

			if !tt.wantErr && !columnsEql(actual, tt.expected) {
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

			if !tt.wantErr && !columnsEql(actual, tt.expected) {
				t.Errorf("expected %v to equal %v", actual, tt.expected)
			}
		})
	}
}

var csv = fakedata.NewColumnFormatter(",")
var def = fakedata.NewColumnFormatter(" ")
var tab = fakedata.NewColumnFormatter("\t")

type args struct {
	input     []string
	formatter fakedata.Formatter
}

func TestGenerateRow(t *testing.T) {
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			"email",
			args{[]string{"email"}, def},
			`^.+?@.+?\..+$`,
		},
		{
			"domain",
			args{[]string{"domain"}, def},
			`^.+?\..+?$`,
		},
		{
			"username",
			args{[]string{"username"}, def},
			`^[a-zA-Z0-9]{2,}$`,
		},
		{
			"double",
			args{[]string{"double"}, def},
			`^-?[0-9]+?(\.[0-9]+?)?$`,
		},
		{
			"date",
			args{[]string{"date"}, def},
			`^\d{4}-\d{2}-\d{2}$`,
		},
		{
			"username domain",
			args{[]string{"username", "domain"}, def},
			`^[a-zA-Z0-9]{2,} .+?\..+?$`,
		},
		{
			"username domain csv",
			args{[]string{"username", "domain"}, csv},
			`^[a-zA-Z0-9]{2,},.+?\..+?$`,
		},
		{
			"username domain tab",
			args{[]string{"username", "domain"}, tab},
			`^[a-zA-Z0-9]{2,}\t.+?\..+?$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns, err := fakedata.NewColumns(tt.args.input)
			if err != nil {
				t.Error(err.Error())
			}
			actual := strings.TrimRight(columns.GenerateRow(tt.args.formatter), "\n")

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
			"int:1,10",
			args{[]string{"int:1,10"}, def},
			1,
			100,
		},
		{
			"int:100,200",
			args{[]string{"int:100,1200"}, def},
			100,
			1200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// this isn't an accurate way of testing random output
			// but it serves a practical purpose
			for index := 0; index < 10000; index++ {
				columns, err := fakedata.NewColumns(tt.args.input)
				if err != nil {
					t.Fatal(err.Error())
				}
				row := strings.Split(columns.GenerateRow(tt.args.formatter), " ")
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
			"date:2016-01-01,2016-12-31",
			args{[]string{"date:2016-01-01,2016-12-31"}, def},
			time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2016, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			"date:2016-01-01,",
			args{[]string{"date:2016-01-01,"}, def},
			time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Now(),
		},
		{
			"date:2046-01-01,2047-01-01",
			args{[]string{"date:2046-01-01,2047-01-01"}, def},
			time.Date(2046, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2047, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// this isn't an accurate way of testing random output
			// but it serves a practical purpose
			for index := 0; index < 10000; index++ {
				columns, err := fakedata.NewColumns(tt.args.input)
				if err != nil {
					t.Fatal(err.Error())
				}

				row := strings.Split(columns.GenerateRow(tt.args.formatter), " ")

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
			args{[]string{"enum"}, def},
			[]string{"foo", "bar", "baz"},
		},
		{
			"enum:Peter,Olivia,Walter",
			args{[]string{"enum:Peter,Olivia,Walter"}, def},
			[]string{"Peter", "Olivia", "Walter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// this isn't an accurate way of testing random output
			// but it serves a practical purpose
			for index := 0; index < 10000; index++ {
				columns, err := fakedata.NewColumns(tt.args.input)
				if err != nil {
					t.Fatal(err.Error())
				}
				row := strings.TrimRight(columns.GenerateRow(tt.args.formatter), "\n")

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
