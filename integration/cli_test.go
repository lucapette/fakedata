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

type testFile struct {
	t    *testing.T
	name string
	dir  string
}

func newFixture(t *testing.T, name string) *testFile {
	return &testFile{t: t, name: name, dir: "fixtures"}
}

func newGoldenFile(t *testing.T, name string) *testFile {
	return &testFile{t: t, name: name, dir: "golden"}
}

func (tf *testFile) path() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		tf.t.Fatal("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
}

func (tf *testFile) write(content string) {
	err := ioutil.WriteFile(tf.path(), []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

func (tf *testFile) asFile() *os.File {
	file, err := os.Open(tf.path())
	if err != nil {
		tf.t.Fatalf("could not open %s: %v", tf.name, err)
	}
	return file
}

func (tf *testFile) load() string {
	content, err := ioutil.ReadFile(tf.path())
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}

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
			"unknown format",
			[]string{"-f=sqll", "-t=USERS", "int:42,42", "enum:foo,foo"},
			"unknown-format.golden",
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
			actual := string(output)

			golden := newGoldenFile(t, tt.golden)

			if *update {
				golden.write(actual)
			}
			expected := golden.load()

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
		golden  string
		wantErr bool
	}{
		{"no file", []string{"file"}, "path-empty.golden", true},
		{"file does not exist", []string{`file:'this file does not exist.txt'`}, "file-does-not-exist.golden", true},
		{"file exists", []string{`file:integration/fixtures/file.txt`}, "file-exist.golden", false},
		{"file exists with quotes", []string{`file:'integration/fixtures/file.txt'`}, "file-exist.golden", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, tt.wantErr, err != nil, err)
			}

			golden := newGoldenFile(t, tt.golden)
			actual := string(output)
			if *update {
				golden.write(actual)
			}

			expected := golden.load()

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		})
	}
}

var templateTests = []struct {
	tmpl    string
	golden  string
	wantErr bool
}{
	{"simple.tmpl", "simple-template.golden", false},
	{"loop.tmpl", "loop.golden", false},
	{"broken.tmpl", "broken-template.golden", true},
	{"unknown-function.tmpl", "unknown-function.golden", true},
}

func TestTemplatesWithCLIArgs(t *testing.T) {
	for _, tt := range templateTests {
		t.Run(tt.tmpl, func(t *testing.T) {
			cmd := exec.Command(binaryPath, "--template", fmt.Sprintf("integration/fixtures/%s", tt.tmpl))
			output, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, tt.wantErr, err != nil, err)
			}

			golden := newGoldenFile(t, tt.golden)
			actual := string(output)
			if *update {
				golden.write(actual)
			}

			expected := golden.load()

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		})
	}
}

func TestTemplatesWithPipe(t *testing.T) {
	for _, tt := range templateTests {
		t.Run(tt.tmpl, func(t *testing.T) {
			fixture := newFixture(t, tt.tmpl)
			cmd := exec.Command(binaryPath)
			cmd.Stdin = fixture.asFile()
			output, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, tt.wantErr, err != nil, err)
			}

			golden := newGoldenFile(t, tt.golden)
			actual := string(output)
			if *update {
				golden.write(actual)
			}

			expected := golden.load()

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

	if err := exec.Command("make").Run(); err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
