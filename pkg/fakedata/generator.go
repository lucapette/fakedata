package fakedata

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

// A Generator is a func that generates random data along with its description
type Generator struct {
	Func func(Column) string
	Desc string
	Name string
}

var generators map[string]Generator

func (g Generator) String() string {
	return fmt.Sprintf("%s\t%s", g.Name, g.Desc)
}

func generate(column Column) string {
	if gen, ok := generators[column.Key]; ok {
		return gen.Func(column)
	}

	return ""
}

// Generators returns all the available generators
func Generators() []Generator {
	gens := make([]Generator, 0)

	for _, v := range generators {
		gens = append(gens, v)
	}

	sort.Slice(gens, func(i, j int) bool { return strings.Compare(gens[i].Name, gens[j].Name) < 0 })
	return gens
}

func withDictKey(key string) func(Column) string {
	return func(column Column) string {
		return dict[key][rand.Intn(len(dict[key]))]
	}
}

func withSep(left, right Column, sep string) func(column Column) string {
	return func(column Column) string {
		return fmt.Sprintf("%s%s%s", generate(left), sep, generate(right))
	}
}

var date = func(column Column) string {
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	if len(column.Min) > 0 {
		if len(column.Max) > 0 {
			formattedMax := fmt.Sprintf("%sT00:00:00.000Z", column.Max)

			date, err := time.Parse("2006-01-02T15:04:05.000Z", formattedMax)
			if err != nil {
				log.Fatalf("Problem with Max: %s", err.Error())
			}

			endDate = date
		}

		formattedMin := fmt.Sprintf("%sT00:00:00.000Z", column.Min)

		date, err := time.Parse("2006-01-02T15:04:05.000Z", formattedMin)
		if err != nil {
			log.Fatal(err.Error())
		}

		startDate = date
	}

	if startDate.After(endDate) {
		log.Fatalf("%v is after %v", startDate, endDate)
	}

	return startDate.Add(time.Duration(rand.Intn(int(endDate.Sub(startDate))))).Format("2006-01-02")
}

var ipv4 = func(column Column) string {
	return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253))
}

var ipv6 = func(column Column) string {
	return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

var mac = func(column Column) string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

var latitute = func(column Column) string {
	lattitude := (rand.Float64() * 180) - 90
	return strconv.FormatFloat(lattitude, 'f', 6, 64)
}

var longitude = func(column Column) string {
	longitude := (rand.Float64() * 360) - 180
	return strconv.FormatFloat(longitude, 'f', 6, 64)
}

var double = func(column Column) string {
	return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64)
}

var integer = func(column Column) string {
	min := 0
	max := 1000

	if len(column.Min) > 0 {
		m, err := strconv.Atoi(column.Min)
		if err != nil {
			log.Fatalf("could not convert min: %v", err)
		}

		min = m

		if len(column.Max) > 0 {
			m, err := strconv.Atoi(column.Max)
			if err != nil {
				log.Fatalf("could not convert max: %v", err)

			}

			max = m
		}
	}

	if min > max {
		log.Fatalf("max(%d) is smaller than min(%d) in %v", max, min, column)
	}

	return strconv.Itoa(min + rand.Intn(max+1-min))
}

func init() {
	generators = make(map[string]Generator)

	generators["date"] = Generator{
		Name: "date",
		Desc: "YYYY-MM-DD. Accepts a range in the format YYYY-MM-DD..YYYY-MM-DD. By default, it generates dates in the last year.",
		Func: date,
	}

	generators["domain.tld"] = Generator{
		Name: "domain.tld",
		Desc: "name|info|com|org|me|us",
		Func: withDictKey("domain.tld"),
	}

	generators["domain.name"] = Generator{
		Name: "domain.tld",
		Desc: "example|test",
		Func: withDictKey("domain.name"),
	}

	generators["country"] = Generator{
		Name: "country",
		Desc: "Full country name",
		Func: withDictKey("country"),
	}

	generators["country.code"] = Generator{
		Name: "country.code",
		Desc: "2-digit country code",
		Func: withDictKey("country.code"),
	}

	generators["state"] = Generator{
		Name: "state",
		Desc: "Full US state name",
		Func: withDictKey("state"),
	}

	generators["state.code"] = Generator{
		Name: "state.code",
		Desc: "2-digit US state name",
		Func: withDictKey("state.code"),
	}

	generators["timezone"] = Generator{
		Name: "timezone",
		Desc: "tz in the form Area/City",
		Func: withDictKey("timezone"),
	}

	generators["username"] = Generator{
		Name: "username",
		Desc: `username using the pattern \w+`,
		Func: withDictKey("username"),
	}

	generators["name.first"] = Generator{
		Name: "name.first",
		Desc: "capilized first name",
		Func: withDictKey("name.first"),
	}

	generators["name.last"] = Generator{
		Name: "name.last",
		Desc: "capilized last name",
		Func: withDictKey("name.last"),
	}

	generators["color"] = Generator{
		Name: "color",
		Desc: "one word color",
		Func: withDictKey("color"),
	}

	generators["product.category"] = Generator{
		Name: "product.category",
		Desc: "Beauty|Games|Movies|Tools|..",
		Func: withDictKey("product.category"),
	}

	generators["product.name"] = Generator{
		Name: "product.name",
		Desc: "invented product name",
		Func: withDictKey("product.name"),
	}

	generators["event.action"] = Generator{
		Name: "event.action",
		Desc: `Clicked|Purchased|Viewed|Watched`,
		Func: withDictKey("event.action"),
	}

	generators["http.method"] = Generator{
		Name: "http.method",
		Desc: `GET|POST|PUT|PATCH|HEAD|DELETE|OPTION`,
		Func: withDictKey("http.method"),
	}

	generators["name"] = Generator{
		Name: "name",
		Desc: `name.first + " " + name.last`,
		Func: withSep(Column{Key: "name.first"}, Column{Key: "name.last"}, " "),
	}

	generators["email"] = Generator{
		Name: "email",
		Desc: "email",
		Func: withSep(Column{Key: "username"}, Column{Key: "domain"}, "@"),
	}

	generators["domain"] = Generator{
		Name: "domain",
		Desc: "domain",
		Func: withSep(Column{Key: "domain.name"}, Column{Key: "domain.tld"}, "."),
	}

	generators["ipv4"] = Generator{Name: "ipv4", Desc: "ipv4", Func: ipv4}

	generators["ipv6"] = Generator{Name: "ipv6", Desc: "ipv6", Func: ipv6}

	generators["mac.address"] = Generator{
		Name: "mac.address",
		Desc: "mac address",
		Func: mac}

	generators["latitute"] = Generator{
		Name: "latitute",
		Desc: "latitute",
		Func: latitute,
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

	generators["int"] = Generator{
		Name: "int",
		Desc: "positive integer. Accepts range mix..max (default: 1..1000).",
		Func: integer,
	}
}
