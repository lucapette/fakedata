package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/lucapette/fakedata/pkg/fakedata"
	flag "github.com/spf13/pflag"
)

var version = "master"

func getFormatter(format, table string) (f fakedata.Formatter, err error) {
	switch format {
	case "csv":
		f = fakedata.NewSeparatorFormatter(",")
	case "tab":
		f = fakedata.NewSeparatorFormatter("\t")
	case "sql":
		f = fakedata.NewSQLFormatter(table)
	case "":
		f = fakedata.NewSeparatorFormatter(" ")
	default:
		err = fmt.Errorf("unknown format: %s", format)
	}
	return f, err
}

func generatorsHelp(generators fakedata.Generators) string {
	max := 0
	for _, gen := range generators {
		if len(gen.Name) > max {
			max = len(gen.Name)
		}
	}

	buffer := &bytes.Buffer{}
	pattern := fmt.Sprintf("%%-%ds%%s\n", max+2) //+2 makes the output more readable
	for _, gen := range generators {
		fmt.Fprintf(buffer, pattern, gen.Name, gen.Desc)
	}

	return buffer.String()
}

func isPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("error checking shell pipe: %v", err)
	}
	// Check if template data is piped to fakedata
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func findTemplate(path string) string {
	if path != "" {
		tp, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("unable to read input: %s", err)
			os.Exit(1)
		}

		return string(tp)
	}

	if isPipe() {
		tp, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("unable to read input: %s", err)
			os.Exit(1)
		}

		return string(tp)
	}

	return ""
}

func main() {
	var (
		generatorsFlag  = flag.BoolP("generators", "G", false, "lists available generators")
		generatorFlag   = flag.StringP("generator", "g", "", "show help for a specific generator")
		constraintsFlag = flag.BoolP("generators-with-constraints", "c", false, "lists available generators with constraints")
		limitFlag       = flag.IntP("limit", "l", 10, "limits rows up to n")
		formatFlag      = flag.StringP("format", "f", "", "generators rows in f format. Available formats: csv|tab|sql")
		versionFlag     = flag.BoolP("version", "v", false, "shows version information")
		tableFlag       = flag.StringP("table", "t", "TABLE", "table name of the sql format")
		templateFlag    = flag.StringP("template", "T", "", "Use template as input")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: fakedata [option ...] field...\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	generators := fakedata.NewGenerators()

	if *generatorsFlag {
		fmt.Print(generatorsHelp(generators))
		os.Exit(0)
	}

	if *generatorFlag != "" {
		if generator := generators.FindByName(*generatorFlag); generator != nil {
			fmt.Printf("Description: %s\n\nExample:\n\n", generator.Desc)
			for i := 0; i < 5; i++ {
				fmt.Println(generator.Func())
			}
		}
		os.Exit(0)
	}

	if *constraintsFlag {
		fmt.Print(generatorsHelp(generators.WithConstraints()))
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	if tmpl := findTemplate(*templateFlag); tmpl != "" {
		if err := fakedata.ExecuteTemplate(tmpl, *limitFlag); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	columns, err := fakedata.NewColumns(flag.Args())
	if err != nil {
		fmt.Printf("%v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	formatter, err := getFormatter(*formatFlag, *tableFlag)
	if err != nil {
		fmt.Printf("%v\n\n", err)
		flag.Usage()
		os.Exit(0)
	}

	for i := 0; i < *limitFlag; i++ {
		fmt.Println(columns.GenerateRow(formatter))
	}
}
