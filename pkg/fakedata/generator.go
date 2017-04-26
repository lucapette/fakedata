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

type generator struct {
	f    func(Column) string
	desc string
}

var generators map[string]generator

func generate(column Column) string {
	if gen, ok := generators[column.Key]; ok {
		return gen.f(column)
	}

	return ""
}

// Generators returns all the available generators
func Generators() []string {
	gens := make([]string, 0)

	for k := range generators {
		gens = append(gens, k)
	}

	sort.Strings(gens)
	return gens
}

func date() func(Column) string {
	return func(column Column) string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
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

func ipv4() func(Column) string {
	return func(column Column) string {
		return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253))
	}

}

func ipv6() func(Column) string {
	return func(column Column) string {
		return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}

}

func mac() func(Column) string {
	return func(column Column) string {
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}
}

func latitute() func(Column) string {
	return func(column Column) string {
		lattitude := (rand.Float64() * 180) - 90
		return strconv.FormatFloat(lattitude, 'f', 6, 64)
	}
}

func longitude() func(Column) string {
	return func(column Column) string {
		longitude := (rand.Float64() * 360) - 180
		return strconv.FormatFloat(longitude, 'f', 6, 64)
	}
}

func double() func(Column) string {
	return func(column Column) string {
		return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64)
	}
}

func integer() func(Column) string {
	return func(column Column) string {
		min := 0
		max := 1000

		if len(column.Range) > 0 {
			rng := strings.Split(column.Range, "..")

			m, err := strconv.Atoi(rng[0])

			if err != nil {
				log.Fatal(err.Error())
			}
			min = m

			if len(rng) > 1 && len(rng[1]) > 0 {
				m, err := strconv.Atoi(rng[1])

				if err != nil {
					log.Fatal(err.Error())
				}
				max = m
			}
		}

		if min > max {
			log.Fatalf("%d is smaller than %d in Column(%s=%s)", max, min, column.Name, column.Key)
		}
		return strconv.Itoa(min + rand.Intn(max-min))
	}
}

func init() {
	generators = make(map[string]generator)

	generators["date"] = generator{desc: "date", f: date()}

	for key := range dict {
		generators[key] = generator{desc: key, f: withDictKey(key)}
	}

	generators["name"] = generator{desc: "name", f: withSep(Column{Key: "name.first"}, Column{Key: "name.last"}, " ")}
	generators["email"] = generator{desc: "email", f: withSep(Column{Key: "username"}, Column{Key: "domain"}, "@")}
	generators["domain"] = generator{desc: "domain", f: withSep(Column{Key: "domain.name"}, Column{Key: "domain.tld"}, ".")}

	generators["ipv4"] = generator{desc: "ipv4", f: ipv4()}
	generators["ipv6"] = generator{desc: "ipv4", f: ipv6()}

	generators["mac.address"] = generator{desc: "mac address", f: mac()}

	generators["latitute"] = generator{desc: "lat", f: latitute()}
	generators["longitude"] = generator{desc: "longitude", f: longitude()}

	generators["double"] = generator{desc: "double", f: double()}

	generators["int"] = generator{desc: "integer generator", f: integer()}
}
