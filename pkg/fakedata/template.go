package fakedata

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

var fakeColumn = Column{"Fake", "fakecolumn", ""}

var generatorFunctions = template.FuncMap{
	"Name": func() string {
		return generators["name"].Func(fakeColumn)
	},
	"Email": func() string {
		return generators["email"].Func(fakeColumn)
	},
	"Int": func(a, b int) string {
		return generators["int"].Func(Column{"Int", "int", fmt.Sprintf("%d..%d", a, b)})
	},
	"Enum": func(keywords ...string) string {
		constraints := createConstraints(keywords)
		return generators["enum"].Func(Column{"Enum", "enu", constraints})
	},
	"Generator": func(name string) string {
		return generators[name].Func(fakeColumn)
	},
	"Loop": func(i int) []int {
		c := make([]int, i)

		return c
	},
	"Odd": func(i int) bool {
		if i%2 != 0 {
			return true
		}
		return false
	},
	"Even": func(i int) bool {
		if i%2 == 0 {
			return true
		}
		return false
	},
}

func createConstraints(params []string) string {
	return strings.Join(params, "..")
}

func getTemplateNameFromPath(name string) string {
	ts := strings.Split(name, "/")
	tn := ts[len(ts)-1]
	return tn
}

func ParseTemplate(path string) {
	tn := getTemplateNameFromPath(path)
	tmp, err := template.New(tn).Funcs(generatorFunctions).ParseFiles(path)
	if err != nil {
		fmt.Println(err)
	}

	tmp.Execute(os.Stdout, generators)

}
