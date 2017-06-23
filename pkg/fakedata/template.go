package fakedata

import (
	"os"
	"strings"
	"text/template"
)

type templateFactory struct {
	factory
}

func newTemplateFactory() *templateFactory {
	return &templateFactory{factory: newFactory()}
}

func (tf templateFactory) getFunctions() template.FuncMap {
	funcMap := template.FuncMap{
		"Loop": func(i int) []int { return make([]int, i) },
		"Odd":  func(i int) bool { return i%2 != 0 },
		"Even": func(i int) bool { return i%2 == 0 },
	}

	for _, gen := range tf.factory.generators {
		if !gen.AcceptsOptions {
			name := strings.Replace(strings.Title(strings.Replace(gen.Name, ".", " ", -1)), " ", "", -1)
			funcMap[name] = gen.Func
		}
	}

	funcMap["Int"] = func(ranges ...int) string {
		min := 0
		max := 1000

		switch len(ranges) {
		case 1:
			min = ranges[0]
		case 2:
			min = ranges[0]
			max = ranges[1]
		}

		return _integer(min, max)
	}

	funcMap["Enum"] = func(enum ...string) string { return withList(enum)() }

	funcMap["File"] = func(path string) (string, error) {
		f, err := file(path)
		if err != nil {
			return "", err
		}

		return f(), nil
	}

	return funcMap
}

func ExecuteTemplate(tmpl string, limit int) (err error) {
	f := newTemplateFactory()
	t, err := template.New("template").Funcs(f.getFunctions()).Parse(tmpl)
	if err != nil {
		return err
	}

	for i := 1; i <= limit; i++ {
		err = t.Execute(os.Stdout, nil)
		if err != nil {
			return err
		}
	}
	return err
}
