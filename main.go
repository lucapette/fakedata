package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/lucapette/fakedata/pkg/fakedata"
	flag "github.com/spf13/pflag"
)

var version = "main"

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
		tp, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("unable to read input: %s", err)
			os.Exit(1)
		}

		return string(tp)
	}

	if isPipe() {
		tp, err := io.ReadAll(os.Stdin)
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
		completionFlag  = flag.StringP("completion", "C", "", "print shell completion function, pass shell name as argument (\"bash\", \"zsh\" or \"fish\")")
		constraintsFlag = flag.BoolP("generators-with-constraints", "c", false, "lists available generators with constraints")
		formatFlag      = flag.StringP("format", "f", "column", "generates rows in f format. Available formats: column|sql")
		generatorFlag   = flag.StringP("generator", "g", "", "show help for a specific generator")
		generatorsFlag  = flag.BoolP("generators", "G", false, "lists available generators")
		headerFlag      = flag.BoolP("header", "H", false, "adds headers row")
		helpFlag        = flag.BoolP("help", "h", false, "shows help")
		limitFlag       = flag.IntP("limit", "l", 10, "limits rows up to n")
		separatorFlag   = flag.StringP("separator", "s", " ", "specifies separator for the column format")
		streamFlag      = flag.BoolP("stream", "S", false, "streams rows till the end of time")
		tableFlag       = flag.StringP("table", "t", "TABLE", "table name of the sql format")
		templateFlag    = flag.StringP("template", "T", "", "Use template as input")
		versionFlag     = flag.BoolP("version", "v", false, "shows version information")
	)

	flag.Usage = func() {
		fmt.Print("Usage: fakedata [option ...] generator...\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *completionFlag != "" {
		completion, err := fakedata.GetCompletionFunc(*completionFlag)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", completion)
		os.Exit(0)
	}

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
				fn := generator.Func
				if generator.IsCustom() {
					custom, err := generator.CustomFunc("")
					if err != nil {
						fmt.Printf("could not generate example: %v", err)
						os.Exit(1)
					}

					fn = custom
				}
				fmt.Println(fn())
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
		if err := fakedata.ExecuteTemplate(tmpl, *limitFlag, *streamFlag); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
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

	var formatter fakedata.Formatter
	if *formatFlag == "column" {
		formatter = fakedata.NewColumnFormatter(*separatorFlag)
	} else if *formatFlag == "sql" {
		formatter = fakedata.NewSQLFormatter(*tableFlag)
	} else {
		fmt.Printf("unknown format: %s\n\n", *formatFlag)
		flag.Usage()
		os.Exit(1)
	}

	fOut := bufio.NewWriter(os.Stdout)
	defer fOut.Flush()

	if *headerFlag {
		columns.GenerateHeader(fOut, formatter)
	}

	if *streamFlag {
		for {
			columns.GenerateRow(fOut, formatter)
		}
	}
	for i := 0; i < *limitFlag; i++ {
		columns.GenerateRow(fOut, formatter)
	}
}
