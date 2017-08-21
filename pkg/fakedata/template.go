package fakedata

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"math/rand"
)

type templateFactory struct {
	factory
}

func newTemplateFactory() *templateFactory {
	return &templateFactory{factory: newFactory()}
}

func (tf templateFactory) getFunctions() template.FuncMap {
	funcMap := template.FuncMap{
		"Loop": func(minmax ...int) []int {
			var size int;

			if len(minmax) == 1 {
				size = minmax[0];
			} else {
				min := minmax[0];
				max := minmax[1];
				size = rand.Intn(max - min) + min;
			}

			return make([]int, size)
		},
		"Odd":  func(i int) bool { return i%2 != 0 },
		"Even": func(i int) bool { return i%2 == 0 },
	}

	for _, gen := range tf.factory.generators {
		name := strings.Replace(strings.Title(strings.Replace(gen.Name, ".", " ", -1)), " ", "", -1)
		if !gen.IsCustom() {
			funcMap[name] = gen.Func
		}
	}

	handler := func(in func(string) (func() string, error), options []string) (string, error) {
		f, err := in(strings.Join(options, ","))
		if err != nil {
			return "", err
		}

		return f(), nil
	}

	funcMap["Int"] = func(ranges ...int) (string, error) {
		options := make([]string, len(ranges))
		for i, r := range ranges {
			options[i] = fmt.Sprintf("%v", r)
		}
		return handler(integer, options)
	}

	funcMap["Enum"] = func(options ...string) (string, error) {
		return handler(enum, options)
	}

	funcMap["File"] = func(path string) (string, error) {
		return handler(file, []string{path})
	}

	funcMap["Date"] = func(dates ...string) (string, error) {
		return handler(date, dates)
	}

	return funcMap
}

// ExecuteTemplate takes a tmpl string and a n int and generates n rows of based
// on the specified tmpl
func ExecuteTemplate(tmpl string, n int) (err error) {
	f := newTemplateFactory()
	t, err := template.New("template").Funcs(f.getFunctions()).Parse(tmpl)
	if err != nil {
		return err
	}

	for i := 1; i <= n; i++ {
		err = t.Execute(os.Stdout, nil)
		if err != nil {
			return err
		}
	}
	return err
}
