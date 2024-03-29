package integration

import (
	"regexp"
	"testing"

	"reflect"

	"github.com/lucapette/fakedata/testutil"
)

func TestCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		golden  string
		wantErr bool
	}{
		{
			"no arguments",
			[]string{},
			"help.golden",
			false,
		},
		{
			"default format",
			[]string{"int:42,42", "enum:foo,foo"},
			"default-format.golden",
			false,
		},
		{
			"default format with headers",
			[]string{"--header", "int:42,42", "enum:foo,foo"},
			"default-format-with-headers.golden",
			false,
		},
		{
			"unknown generators",
			[]string{"madeupgenerator", "anothermadeupgenerator"},
			"unknown-generators.golden",
			true,
		},
		{
			"default format with limit short",
			[]string{"-l=5", "int:42,42", "enum:foo,foo"},
			"default-format-with-limit.golden",
			false,
		},
		{
			"default format with limit",
			[]string{"--limit=5", "int:42,42", "enum:foo,foo"},
			"default-format-with-limit.golden",
			false,
		},
		{
			"csv format short",
			[]string{"-s=,", "int:42,42", "enum:foo,foo"},
			"csv-format.golden",
			false,
		},
		{
			"csv format",
			[]string{"--separator=,", "int:42,42", "enum:foo,foo"},
			"csv-format.golden",
			false,
		},
		{
			"tab format",
			[]string{"--separator=\t", "int:42,42", "enum:foo,foo"},
			"tab-format.golden",
			false,
		},
		{
			"sql format",
			[]string{"-f=sql", "int:42,42", "enum:foo,foo"},
			"sql-format.golden",
			false,
		},
		{
			"sql format with keys",
			[]string{"-f=sql", "age=int:42,42", "name=enum:foo,foo"},
			"sql-format-with-keys.golden",
			false,
		},
		{
			"sql format with table name",
			[]string{"-f=sql", "-t=USERS", "int:42,42", "enum:foo,foo"},
			"sql-format-with-table-name.golden",
			false,
		},
		{
			"unknown format",
			[]string{"-f=no-format", "-t=USERS", "int:42,42", "enum:foo,foo"},
			"unknown-format.golden",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, tt.wantErr, err != nil, err)
			}
			actual := string(output)

			golden := testutil.NewGoldenFile(t, tt.golden)

			if *update {
				golden.Write(actual)
			}
			expected := golden.Load()

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}

func TestGeneratorDescription(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"simple generator", []string{"-g", "name.first"}},
		{"custom generator", []string{"-g", "int"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if err != nil {
				t.Fatalf("test run returned an error: %v\n%s", err, output)
			}

			actual := string(output)
			matched, err := regexp.MatchString("Description:", actual)
			if err != nil {
				t.Fatalf("could not match actual: %v", err)
			}

			if !matched {
				t.Fatalf("expected %s to match description, but did not", actual)
			}
		})
	}
}

func TestFileGenerator(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		golden  string
		wantErr bool
	}{
		{"no file", []string{"file"}, "path-empty.golden", true},
		{"file does not exist", []string{`file:'this file does not exist.txt'`}, "file-does-not-exist.golden", true},
		{"file exists", []string{`file:testutil/fixtures/file.txt`}, "file-exist.golden", false},
		{"file exists with quotes", []string{`file:'testutil/fixtures/file.txt'`}, "file-exist.golden", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runBinary(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, tt.wantErr, err != nil, err)
			}

			golden := testutil.NewGoldenFile(t, tt.golden)
			actual := string(output)
			if *update {
				golden.Write(actual)
			}

			expected := golden.Load()

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}
