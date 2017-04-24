package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/lucapette/testdata/pkg/testdata"
	flag "github.com/spf13/pflag"
)

var usage = `
  Usage: testdata
    [--tick d]
    [--max n]
    [--list]

    testdata -h | --help

  Options:
    --list          list all available generators
    --max n         generate data up to n [default: 10]
    --tick d        generate data every d [default: 10ms]
    -v, --version   show version information
    -h, --help      show help information
`

var listFlag = flag.Bool("list", false, "list all the generators")
var tickFlag = flag.Duration("tick", 10*time.Millisecond, "generate data every d milliseconds")
var maxFlag = flag.Int("max", 10, "generate up to n rows")
var helpFlag = flag.Bool("help", false, "print usage")
var formatFlag = flag.String("format", "", "Output format")

func main() {
	if *helpFlag {
		fmt.Printf(usage)
		os.Exit(0)
	}

	if *listFlag {
		generators := testdata.List()
		sort.Strings(generators)

		for _, name := range generators {
			fmt.Printf("%s\n", name)
		}
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	tick := time.Tick(*tickFlag)

	total := 0

	for _ = range tick {
		fmt.Print(testdata.GenerateRow(flag.Args(), *formatFlag))

		if total++; total == *maxFlag {
			return
		}
	}
}

func init() {
	flag.Parse()
}
