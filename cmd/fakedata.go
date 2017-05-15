package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lucapette/fakedata/pkg/fakedata"
	flag "github.com/spf13/pflag"
)

var version = "master"

var generatorsFlag = flag.BoolP("generators", "g", false, "lists available generators")
var limitFlag = flag.IntP("limit", "l", 10, "limits rows up to n")
var formatFlag = flag.StringP("format", "f", "", "generators rows in f format. Available formats: csv|tab|sql")
var versionFlag = flag.BoolP("version", "v", false, "shows version information")
var tableFlag = flag.StringP("table", "t", "TABLE", "table name of the sql format")

func getFormatter(format string) (f fakedata.Formatter) {
	switch format {
	case "csv":
		f = fakedata.NewSeparatorFormatter(",")
	case "tab":
		f = fakedata.NewSeparatorFormatter("\t")
	case "sql":
		f = fakedata.NewSQLFormatter(*tableFlag)
	default:
		f = fakedata.NewSeparatorFormatter(" ")
	}
	return f
}

func generatorsHelp(generators []fakedata.Generator) string {
	var buffer bytes.Buffer

	var max int

	for _, gen := range generators {
		if len(gen.Name) > max {
			max = len(gen.Name)
		}
	}

	pattern := fmt.Sprintf("%%-%ds%%s\n", max+2) //+2 makes the output more readable
	for _, gen := range generators {
		buffer.WriteString(fmt.Sprintf(pattern, gen.Name, gen.Desc))
	}

	return buffer.String()
}

func main() {
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if *generatorsFlag {
		fmt.Print(generatorsHelp(fakedata.Generators()))
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if err := fakedata.ValidateGenerators(flag.Args()); err != nil {
		fmt.Printf("%v\n\nSee fakedata --generators for a list of available generators.\n", err)
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	columns := fakedata.NewColumns(flag.Args())
	formatter := getFormatter(*formatFlag)

	for i := 0; i < *limitFlag; i++ {
		fmt.Print(fakedata.GenerateRow(columns, formatter))
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [option ...] field...\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}
