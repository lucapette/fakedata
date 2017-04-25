package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lucapette/fakedata/pkg/fakedata"
	flag "github.com/spf13/pflag"
)

var version = "master"

var usage = `
  Usage: fakedata [option ...] [field ...]

  Options:
    --generators    lists available generators
    --limit n       limits rows up to n [default: 10]
    --help          shows help information
    --format f      generates rows in f format [options: csv|tab|sql, default: " "]
    --version       shows version information
`

var generatorsFlag = flag.Bool("generators", false, "lists available generators")
var limitFlag = flag.Int("limit", 10, "limits rows up to n")
var helpFlag = flag.Bool("help", false, "shows help information")
var formatFlag = flag.String("format", "", "generators rows in f format")
var versionFlag = flag.Bool("version", false, "shows version information")

func main() {
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if *helpFlag {
		fmt.Print(usage)
		os.Exit(0)
	}

	if *generatorsFlag {
		for _, generator := range fakedata.Generators() {
			fmt.Printf("%s\n", generator)
		}
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Printf(usage)
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	columns := fakedata.NewColumns(flag.Args())
	for i := 0; i < *limitFlag; i++ {
		fmt.Print(fakedata.GenerateRow(columns, *formatFlag))
	}
}

func init() {
	flag.Parse()
}
