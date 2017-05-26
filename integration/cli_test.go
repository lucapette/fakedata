package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"reflect"

	"github.com/kr/pretty"
)

// In the following tests, there's a lot going on.
// Please have a look at the following article for a longer explanation:
// http://lucapette.me/writing-integration-tests-for-a-go-cli-application

var update = flag.Bool("update", false, "update golden files")

const binaryName = "fakedata"

var binaryPath string

func diff(expected, actual interface{}) []string {
	return pretty.Diff(expected, actual)
}

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), fixture)
}

func writeFixture(t *testing.T, fixture string, content []byte) {
	err := ioutil.WriteFile(fixturePath(t, fixture), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func TestCliArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		{
			name:    "no arguments",
			args:    []string{},
			fixture: "help.golden",
		},
		{
			name:    "list generators",
			args:    []string{"-g"},
			fixture: "generators.golden",
		},
		{
			name:    "default format",
			args:    []string{"int,42..42", "enum,foo..foo"},
			fixture: "default-format.golden",
		},
		{
			name:    "unknown generators",
			args:    []string{"madeupgenerator", "anothermadeupgenerator"},
			fixture: "unknown-generators.golden",
		},
		{
			name:    "default format with limit short",
			args:    []string{"-l=5", "int,42..42", "enum,foo..foo"},
			fixture: "default-format-with-limit.golden",
		},
		{
			name:    "default format with limit",
			args:    []string{"--limit=5", "int,42..42", "enum,foo..foo"},
			fixture: "default-format-with-limit.golden",
		},
		{
			name:    "csv format short",
			args:    []string{"-f=csv", "int,42..42", "enum,foo..foo"},
			fixture: "csv-format.golden",
		},
		{
			name:    "csv format",
			args:    []string{"--format=csv", "int,42..42", "enum,foo..foo"},
			fixture: "csv-format.golden",
		},
		{
			name:    "tab format",
			args:    []string{"-f=tab", "int,42..42", "enum,foo..foo"},
			fixture: "tab-format.golden",
		},
		{
			name:    "sql format",
			args:    []string{"-f=sql", "int,42..42", "enum,foo..foo"},
			fixture: "sql-format.golden",
		},
		{
			name:    "sql format with keys",
			args:    []string{"-f=sql", "age=int,42..42", "name=enum,foo..foo"},
			fixture: "sql-format-with-keys.golden",
		},
		{
			name:    "sql format with table name",
			args:    []string{"-f=sql", "-t=USERS", "int,42..42", "enum,foo..foo"},
			fixture: "sql-format-with-table-name.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("output: %s\nerr: %v", output, err)
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)
			expected := loadFixture(t, tt.fixture)

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		})
	}
}

func TestFileGenerator(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
		wantErr bool
	}{
		{"no file", []string{"file"}, "path-empty.golden", true},
		{"no file,", []string{"file,"}, "path-empty.golden", true},
		{"file does not exist", []string{`file,'this file does not exist.txt'`}, "file-empty.golden", true},
		{"file exists", []string{`file,integration/file.txt`}, "file-exist.golden", false},
		{"file exists with quotes", []string{`file,'integration/file.txt'`}, "file-exist.golden", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, tt.wantErr, err != nil, err)
			}
			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)
			expected := loadFixture(t, tt.fixture)

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		})
	}
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	abs, err := filepath.Abs(binaryName)
	if err != nil {
		fmt.Printf("could not get abs path for %s: %v", binaryName, err)
		os.Exit(1)
	}

	binaryPath = abs

	make := exec.Command("make")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
