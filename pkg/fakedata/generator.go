package fakedata

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lucapette/fakedata/pkg/data"
)

// A Generator is a func that generates random data along with its description
type Generator struct {
	Func       func() string
	CustomFunc func(string) (func() string, error)
	Desc       string
	Name       string
}

// IsCustom returns a bool indicating whether the generator has a CustomFunc or
// not
func (g Generator) IsCustom() bool {
	return g.CustomFunc != nil
}

// Generators returns available generators
func Generators() []Generator {
	f := newFactory()
	gens := make([]Generator, 0)

	for _, v := range f.generators {
		gens = append(gens, v)
	}

	sort.Slice(gens, func(i, j int) bool { return strings.Compare(gens[i].Name, gens[j].Name) < 0 })
	return gens
}

func withList(list []string) func() string {
	return func() string {
		return list[rand.Intn(len(list))]
	}
}

func ipv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253))
}

func ipv6() string {
	return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func mac() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func latitude() string {
	return strconv.FormatFloat((rand.Float64()*180)-90, 'f', 6, 64)
}

func longitude() string {
	return strconv.FormatFloat((rand.Float64()*360)-180, 'f', 6, 64)
}

func double() string {
	return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64)
}

func date(options string) (f func() string, err error) {
	var min, max string

	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	dateRange := strings.Split(options, ",")
	min = dateRange[0]

	if len(dateRange) > 1 {
		max = dateRange[1]
	}

	if len(min) > 0 {
		if len(max) > 0 {
			formattedMax := fmt.Sprintf("%sT00:00:00.000Z", max)

			date, err := time.Parse("2006-01-02T15:04:05.000Z", formattedMax)
			if err != nil {
				return nil, fmt.Errorf("problem parsing max date: %v", err)
			}

			endDate = date
		}

		formattedMin := fmt.Sprintf("%sT00:00:00.000Z", min)

		date, err := time.Parse("2006-01-02T15:04:05.000Z", formattedMin)
		if err != nil {
			return nil, fmt.Errorf("problem parsing mix date: %v", err)
		}

		startDate = date
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("%v is after %v", startDate, endDate)
	}

	return func() string {
		return startDate.Add(time.Duration(rand.Intn(int(endDate.Sub(startDate))))).Format("2006-01-02")
	}, err
}

func integer(options string) (func() string, error) {
	min := 0
	max := 1000
	var low, high string
	intRange := strings.Split(options, ",")
	low = intRange[0]

	if len(intRange) > 1 {
		high = intRange[1]
	}

	if len(low) > 0 {
		m, err := strconv.Atoi(low)
		if err != nil {
			return nil, fmt.Errorf("could not convert min: %v", err)
		}

		min = m

		if len(high) > 0 {
			m, err := strconv.Atoi(high)
			if err != nil {
				return nil, fmt.Errorf("could not convert max: %v", err)
			}

			max = m
		}
	}

	if min > max {
		return nil, fmt.Errorf("max(%d) is smaller than min(%d)", max, min)
	}

	return func() string { return strconv.Itoa(min + rand.Intn(max+1-min)) }, nil
}

func file(path string) (func() string, error) {
	if path == "" {
		return nil, fmt.Errorf("no file path given")
	}

	filePath := strings.Trim(path, "'\"")

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", filePath, err)
	}
	list := strings.Split(string(content), "\n")

	return func() string { return withList(list)() }, nil
}

func enum(options string) (func() string, error) {
	list := []string{"foo", "bar", "baz"}
	if options != "" {
		list = strings.Split(options, ",")
	}
	return func() string { return withList(list)() }, nil
}

type factory struct {
	generators map[string]Generator
}

func (f factory) extractFunc(key, options string) (fn func() string, err error) {
	gen, ok := f.generators[key]
	if !ok {
		return nil, fmt.Errorf("unknown generator: %s", key)
	}

	if gen.IsCustom() {
		return gen.CustomFunc(options)
	}

	return gen.Func, nil
}

func domain() string {
	return withList([]string{"test", "example"})() + "." + withList(data.TLDs)()
}

func newFactory() (f factory) {
	generators := make(map[string]Generator)

	generators["domain.tld"] = Generator{
		Name: "domain.tld",
		Desc: "name|info|com|org|me|us",
		Func: withList(data.TLDs),
	}

	generators["domain.name"] = Generator{
		Name: "domain.name",
		Desc: "example|test",
		Func: withList([]string{"example", "test"}),
	}

	generators["country"] = Generator{
		Name: "country",
		Desc: "Full country name",
		Func: withList(data.Countries),
	}

	generators["country.code"] = Generator{
		Name: "country.code",
		Desc: "2-digit country code",
		Func: withList(data.CountryCodes),
	}

	generators["state"] = Generator{
		Name: "state",
		Desc: "Full US state name",
		Func: withList(data.States),
	}

	generators["state.code"] = Generator{
		Name: "state.code",
		Desc: "2-digit US state name",
		Func: withList(data.StateCodes),
	}

	generators["timezone"] = Generator{
		Name: "timezone",
		Desc: "tz in the form Area/City",
		Func: withList(data.Timezones),
	}

	generators["username"] = Generator{
		Name: "username",
		Desc: `username using the pattern \w+`,
		Func: withList(data.Usernames),
	}

	generators["name.first"] = Generator{
		Name: "name.first",
		Desc: "capitalized first name",
		Func: withList(data.Firstnames),
	}

	generators["name.last"] = Generator{
		Name: "name.last",
		Desc: "capitalized last name",
		Func: withList(data.Lastnames),
	}

	generators["color"] = Generator{
		Name: "color",
		Desc: "one word color",
		Func: withList(data.Colors),
	}

	generators["product.category"] = Generator{
		Name: "product.category",
		Desc: "Beauty|Games|Movies|Tools|..",
		Func: withList(data.ProductCategories),
	}

	generators["product.name"] = Generator{
		Name: "product.name",
		Desc: "invented product name",
		Func: withList(data.ProductNames),
	}

	generators["event.action"] = Generator{
		Name: "event.action",
		Desc: `clicked|purchased|viewed|watched`,
		Func: withList([]string{"clicked", "purchased", "viewed", "watched"}),
	}

	generators["http.method"] = Generator{
		Name: "http.method",
		Desc: `DELETE|GET|HEAD|OPTION|PATCH|POST|PUT`,
		Func: withList([]string{"DELETE", "GET", "HEAD", "OPTION", "PATCH", "POST", "PUT"}),
	}

	generators["name"] = Generator{
		Name: "name",
		Desc: `name.first + " " + name.last`,
		Func: func() string {
			return withList(data.Firstnames)() + " " + withList(data.Lastnames)()
		},
	}

	generators["email"] = Generator{
		Name: "email",
		Desc: "email",
		Func: func() string {
			return withList(data.Usernames)() + "@" + domain()
		},
	}

	generators["domain"] = Generator{
		Name: "domain",
		Desc: "domain",
		Func: domain,
	}

	generators["ipv4"] = Generator{Name: "ipv4", Desc: "ipv4", Func: ipv4}

	generators["ipv6"] = Generator{Name: "ipv6", Desc: "ipv6", Func: ipv6}

	generators["mac.address"] = Generator{
		Name: "mac.address",
		Desc: "mac address",
		Func: mac,
	}

	generators["latitude"] = Generator{
		Name: "latitude",
		Desc: "latitude",
		Func: latitude,
	}

	generators["longitude"] = Generator{
		Name: "longitude",
		Desc: "longitude",
		Func: longitude,
	}

	generators["double"] = Generator{
		Name: "double",
		Desc: "double number",
		Func: double,
	}

	generators["noun"] = Generator{
		Name: "noun",
		Desc: "random noun",
		Func: withList(data.Nouns),
	}

	generators["emoji"] = Generator{
		Name: "emoji",
		Desc: "random emoji",
		Func: withList(data.Emoji),
	}

	generators["animal"] = Generator{
		Name: "animal",
		Desc: "random animal name",
		Func: withList(data.Animals),
	}

	generators["animal.cat"] = Generator{
		Name: "animal.cat",
		Desc: "random cat breed",
		Func: withList(data.Cats),
	}

	// custom generators
	generators["date"] = Generator{
		Name:       "date",
		Desc:       "YYYY-MM-DD. Accepts a range in the format YYYY-MM-DD,YYYY-MM-DD. By default, it generates dates in the last year.",
		CustomFunc: date,
	}

	generators["int"] = Generator{
		Name:       "int",
		Desc:       "positive integer. Accepts range min..max (default: 1,1000).",
		CustomFunc: integer,
	}

	generators["enum"] = Generator{
		Name:       "enum",
		Desc:       `a random value from an enum. Defaults to "foo,bar,baz"`,
		CustomFunc: enum,
	}

	generators["file"] = Generator{
		Name:       "file",
		Desc:       `Read a random line from a file. Pass filepath with 'file,path/to/file.txt'.`,
		CustomFunc: file,
	}

	return factory{generators: generators}
}
