package fakedata

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var generators map[string]func(Column) string

func generate(column Column) string {
	if f, ok := generators[column.Key]; ok {
		return f(column)
	}

	return ""
}

func init() {
	generators = make(map[string]func(Column) string)

	withDictKey := func(key string) func(Column) string {
		return func(column Column) string {
			return dict[key][rand.Intn(len(dict[key]))]
		}
	}

	for key := range dict {
		generators[key] = withDictKey(key)
	}

	withSep := func(left, right Column, sep string) func(column Column) string {
		return func(column Column) string {
			return fmt.Sprintf("%s%s%s", generate(left), sep, generate(right))
		}
	}
	generators["name"] = withSep(Column{Key: "name.first"}, Column{Key: "name.last"}, " ")
	generators["email"] = withSep(Column{Key: "username"}, Column{Key: "domain"}, "@")
	generators["domain"] = withSep(Column{Key: "domain.name"}, Column{Key: "domain.tld"}, ".")

	generators["unixtime"] = func(column Column) string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	generators["id"] = func(column Column) string {
		chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		ret := make([]rune, 10)

		for i := range ret {
			ret[i] = chars[rand.Intn(len(chars))]
		}

		return string(ret)
	}

	generators["ipv4"] = func(column Column) string {
		return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253))
	}

	generators["ipv6"] = func(column Column) string {
		return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}

	generators["mac.address"] = func(column Column) string {
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}

	generators["latitude"] = func(column Column) string {
		lattitude := (rand.Float64() * 180) - 90
		return strconv.FormatFloat(lattitude, 'f', 6, 64)
	}

	generators["longitude"] = func(column Column) string {
		longitude := (rand.Float64() * 360) - 180
		return strconv.FormatFloat(longitude, 'f', 6, 64)
	}

	generators["double"] = func(column Column) string {
		return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64)
	}

	generators["int"] = func(column Column) string {
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
