package integration

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// In these tests, there's a lot going on. See this blog post for more details:
// https://lucapette.me/writing/writing-integration-tests-for-a-go-cli-application

var update = flag.Bool("update", false, "update golden files")

const binaryName = "fakedata-with-cover"

var binaryPath string

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

	os.Exit(m.Run())
}

func runBinary(args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}
