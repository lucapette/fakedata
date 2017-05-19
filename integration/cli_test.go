package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"reflect"

	"github.com/kr/pretty"
)

var update = flag.Bool("update", false, "update golden files")

var binaryName = "fakedata"

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
		{"no arguments", []string{}, "help.golden"},
		{"list generators", []string{"-g"}, "generators.golden"},
		{"default format", []string{"int,42..42", "enum,foo..foo"}, "default-format.golden"},
		{"unknown generators", []string{"madeupgenerator", "anothermadeupgenerator"}, "unknown-generators.golden"},
		{"default format with limit short", []string{"-l=5", "int,42..42", "enum,foo..foo"}, "default-format-with-limit.golden"},
		{"default format with limit", []string{"--limit=5", "int,42..42", "enum,foo..foo"}, "default-format-with-limit.golden"},
		{"csv format short", []string{"-f=csv", "int,42..42", "enum,foo..foo"}, "csv-format.golden"},
		{"csv format", []string{"--format=csv", "int,42..42", "enum,foo..foo"}, "csv-format.golden"},
		{"tab format", []string{"-f=tab", "int,42..42", "enum,foo..foo"}, "tab-format.golden"},
		{"sql format", []string{"-f=sql", "int,42..42", "enum,foo..foo"}, "sql-format.golden"},
		{"sql format with keys", []string{"-f=sql", "age=int,42..42", "name=enum,foo..foo"}, "sql-format-with-keys.golden"},
		{"sql format with table name", []string{"-f=sql", "-t=USERS", "int,42..42", "enum,foo..foo"}, "sql-format-with-table-name.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(path.Join(dir, binaryName), tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
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
	make := exec.Command("make")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
