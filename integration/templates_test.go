package integration

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"

	"github.com/lucapette/fakedata/testutil"
)

var templateTests = []struct {
	tmpl    string
	golden  string
	wantErr bool
}{
	{"simple.tmpl", "simple-template.golden", false},
	{"loop.tmpl", "loop.golden", false},
	{"loop-with-index.tmpl", "loop-with-index.golden", false},
	{"broken.tmpl", "broken-template.golden", true},
	{"unknown-function.tmpl", "unknown-function.golden", true},
}

func TestTemplatesWithCLIArgs(t *testing.T) {
	for _, tt := range templateTests {
		t.Run(tt.tmpl, func(t *testing.T) {
			cmd := exec.Command(binaryPath, "--template", fmt.Sprintf("testutil/fixtures/%s", tt.tmpl))
			output, err := cmd.CombinedOutput()
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

func TestTemplatesWithPipe(t *testing.T) {
	for _, tt := range templateTests {
		t.Run(tt.tmpl, func(t *testing.T) {
			fixture := testutil.NewFixture(t, tt.tmpl)
			cmd := exec.Command(binaryPath)
			cmd.Stdin = fixture.AsFile()
			output, err := cmd.CombinedOutput()
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
