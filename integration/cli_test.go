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
	"strings"

	"github.com/kr/pretty"
)

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
		match   func(string, string) bool
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
		{
			name:    "file generator",
			args:    []string{"file,integration/file.golden"},
			fixture: "file.golden",
			match: func(actual, expected string) bool {
				for _, line := range strings.Split(actual, "\n") {
					if !strings.Contains(expected+"\n", line+"\n") {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err, "\n", string(output))
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)
			expected := loadFixture(t, tt.fixture)

			if tt.match != nil {
				if !tt.match(actual, expected) {
					t.Fatalf("values do not match: \n%v\n%v", actual, expected)
				}
			} else if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		})
	}
}

func TestCliErr(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no file", []string{"file"}},
		{"no file,", []string{"file,"}},
		{"no file,''", []string{"file,''"}},
		{`no file,""`, []string{`file,""`}},
		{"file does not exist", []string{`file,'this file does not exist.txt'`}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := exec.Command(binaryPath, tt.args...).Run(); err == nil {
				t.Fatalf("expected to fail with args: %v", tt.args)
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
	make := exec.Command("make")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}
	binaryPath, err = filepath.Abs(binaryName)
	if err != nil {
		fmt.Printf("could not get binary path: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
