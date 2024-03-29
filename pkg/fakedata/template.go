package fakedata

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
			var n int

			if len(minmax) == 1 {
				n = minmax[0]
			} else {
				min := minmax[0]
				max := minmax[1]
				if min == max {
					n = min
				} else {
					n = rand.Intn(max-min) + min
				}
			}

			times := make([]int, n)
			for i := 0; i < n; i++ {
				times[i] = i
			}

			return times
		},
		"Odd":  func(i int) bool { return i%2 != 0 },
		"Even": func(i int) bool { return i%2 == 0 },
	}

	c := cases.Title(language.English)

	for _, gen := range tf.factory.generators {
		if !gen.IsCustom() {
			name := strings.Replace(c.String(strings.Replace(gen.Name, ".", " ", -1)), " ", "", -1)
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
// on the specified tmpl. Will loop forever if streamMode is true
func ExecuteTemplate(tmpl string, n int, streamMode bool) (err error) {
	fOut := bufio.NewWriter(os.Stdout)
	defer fOut.Flush()

	f := newTemplateFactory()
	t, err := template.New("template").Funcs(f.getFunctions()).Parse(tmpl)
	if err != nil {
		return err
	}

	if streamMode {
		for {
			err = t.Execute(fOut, nil)
			if err != nil {
				return err
			}
		}
	}

	for i := 1; i <= n; i++ {
		err = t.Execute(fOut, nil)
		if err != nil {
			return err
		}
	}
	return err
}
