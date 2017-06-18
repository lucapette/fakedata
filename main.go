package main

import (
	"bytes"
	"fmt"
	"github.com/lucapette/fakedata/pkg/fakedata"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"math/rand"
	"os"
	"text/template"
	"time"
)

var version = "master"
var tmp *template.Template

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

func generatorsHelp() string {
	generators := fakedata.Generators()
	var buffer bytes.Buffer

	var max int

	for _, gen := range generators {
		if len(gen.Name) > max {
			max = len(gen.Name)
		}
	}

	pattern := fmt.Sprintf("%%-%ds%%s\n", max+2) //+2 makes the output more readable
	for _, gen := range generators {
		fmt.Fprintf(&buffer, pattern, gen.Name, gen.Desc)
	}

	return buffer.String()
}

func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	// Check if template data is piped to fakedata
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func main() {
	var (
		generatorsFlag = flag.BoolP("generators", "g", false, "lists available generators")
		limitFlag      = flag.IntP("limit", "l", 10, "limits rows up to n")
		formatFlag     = flag.StringP("format", "f", "", "generators rows in f format. Available formats: csv|tab|sql")
		versionFlag    = flag.BoolP("version", "v", false, "shows version information")
		tableFlag      = flag.StringP("table", "t", "TABLE", "table name of the sql format")
		templateFlag   = flag.StringP("template", "", "", "Use template as input")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fakedata [option ...] field...\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if *generatorsFlag {
		fmt.Print(generatorsHelp())
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	if *templateFlag != "" {
		t, err := fakedata.ParseTemplate(*templateFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tmp = t
	}

	if isPipe() {
		tp, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			fmt.Printf("Unable to read input: %s", err)
			os.Exit(1)
		}

		t, err := fakedata.ParseTemplateFromPipe(string(tp))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tmp = t
	}
	// Execute a template if there is one
	if tmp != nil {
		if err := fakedata.ExecuteTemplate(tmp, *limitFlag); err != nil {
			fmt.Println(err)
		}

		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	columns, err := fakedata.NewColumns(flag.Args())
	if err != nil {
		fmt.Printf("%v\n\nSee fakedata --generators for a list of available generators.\n", err)
		os.Exit(0)
	}

	formatter, err := getFormatter(*formatFlag, *tableFlag)
	if err != nil {
		fmt.Printf("%v\n\n", err)
		flag.Usage()
		os.Exit(0)
	}

	for i := 0; i < *limitFlag; i++ {
		fmt.Print(fakedata.GenerateRow(columns, formatter))
	}
}
