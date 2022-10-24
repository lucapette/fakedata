package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type TestFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewFixture(t *testing.T, name string) *TestFile {
	return &TestFile{t: t, name: name, dir: "fixtures"}
}

func NewGoldenFile(t *testing.T, name string) *TestFile {
	return &TestFile{t: t, name: name, dir: "golden"}
}

func (tf *TestFile) Path() string {
	tf.t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		tf.t.Fatal("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
}

func (tf *TestFile) Write(content string) {
	tf.t.Helper()
	err := os.WriteFile(tf.Path(), []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

func (tf *TestFile) AsFile() *os.File {
	tf.t.Helper()
	file, err := os.Open(tf.Path())
	if err != nil {
		tf.t.Fatalf("could not open %s: %v", tf.name, err)
	}
	return file
}

func (tf *TestFile) Load() string {
	tf.t.Helper()

	content, err := os.ReadFile(tf.Path())
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}
