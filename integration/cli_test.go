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

func TestCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
		wantErr bool
	}{
		{
			"no arguments",
			[]string{},
			"help.golden",
			false,
		},
		{
			"list generators",
			[]string{"-g"},
			"generators.golden",
			false,
		},
		{
			"default format",
			[]string{"int:42,42", "enum:foo,foo"},
			"default-format.golden",
			false,
		},
		{
			"unknown generators",
			[]string{"madeupgenerator"},
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
			[]string{"-f=csv", "int:42,42", "enum:foo,foo"},
			"csv-format.golden",
			false,
		},
		{
			"csv format",
			[]string{"--format=csv", "int:42,42", "enum:foo,foo"},
			"csv-format.golden",
			false,
		},
		{
			"tab format",
			[]string{"-f=tab", "int:42,42", "enum:foo,foo"},
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
			"no file",
			[]string{"file"},
			"path-empty.golden",
			true,
		},
		{
			"no file,",
			[]string{"file:"},
			"path-empty.golden",
			true,
		},
		{
			"file does not exist",
			[]string{`file:'this file does not exist.txt'`},
			"file-does-not-exist.golden",
			true,
		},
		{
			"file exists",
			[]string{`file:integration/file.txt`},
			"file-exist.golden",
			false,
		},
		{
			"file exists with quotes",
			[]string{`file:'integration/file.txt'`},
			"file-exist.golden",
			false,
		},
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
