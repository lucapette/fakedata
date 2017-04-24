package fakedata

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var generators map[string]func() string

func init() {
	generators = make(map[string]func() string)

	withDictKey := func(key string) func() string {
		return func() string {
			return dict[key][rand.Intn(len(dict[key]))]
		}
	}

	for key := range dict {
		generators[key] = withDictKey(key)
	}

	withSep := func(left, right, sep string) func() string {
		return func() string {
			return fmt.Sprintf("%s%s%s", Generate(left), sep, Generate(right))
		}
	}
	generators["name"] = withSep("name.first", "name.last", " ")
	generators["email"] = withSep("username", "domain", "@")
	generators["domain"] = withSep("domain.name", "domain.tld", ".")

	generators["unixtime"] = func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	generators["id"] = func() string {
		chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		ret := make([]rune, 10)

		for i := range ret {
			ret[i] = chars[rand.Intn(len(chars))]
		}

		return string(ret)
	}

	generators["ipv4"] = func() string {
		return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(253), rand.Intn(255), rand.Intn(255), 1+rand.Intn(253))
	}

	generators["ipv6"] = func() string {
		return fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}

	generators["mac.address"] = func() string {
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}

	generators["latitude"] = func() string {
		lattitude := (rand.Float64() * 180) - 90
		return strconv.FormatFloat(lattitude, 'f', 6, 64)
	}

	generators["longitude"] = func() string {
		longitude := (rand.Float64() * 360) - 180
		return strconv.FormatFloat(longitude, 'f', 6, 64)
	}

	generators["double"] = func() string {
		return strconv.FormatFloat(rand.NormFloat64()*1000, 'f', 4, 64)
	}
}
