package fakedata

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// generatorFunctions holds all the functions available for the template
var generatorFunctions = template.FuncMap{
	"Loop": func(i int) []int {
		return make([]int, i)
	},
	"Odd": func(i int) bool {
		return i%2 != 0
	},
	"Even": func(i int) bool {
		return i%2 == 0
	},
	"Date": func() string {
		return generators["date"].Func(Column{Name: "tmplDate", Key: "tmplDateKey", Constraints: ""})
	},
	"DomainTld": func() string {
		return generators["domain.tld"].Func(Column{Name: "tmplDomainTld", Key: "tmplDomainTld", Constraints: ""})
	},
	"DomainName": func() string {
		return generators["domain.name"].Func(Column{Name: "tmplDomainName", Key: "tmplDomainName", Constraints: ""})
	},
	"Country": func() string {
		return generators["country"].Func(Column{Name: "tmplCountry", Key: "tmplCountryKey", Constraints: ""})
	},
	"CountryCode": func() string {
		return generators["country.code"].Func(Column{Name: "tmplCountryCode", Key: "tmplCountryCode", Constraints: ""})
	},
	"State": func() string {
		return generators["state"].Func(Column{Name: "tmplState", Key: "tmplStateKey", Constraints: ""})
	},
	"Timezone": func() string {
		return generators["timezone"].Func(Column{Name: "tmplTimezone", Key: "tmplTimezoneKey", Constraints: ""})
	},
	"Username": func() string {
		return generators["username"].Func(Column{Name: "tmplUsername", Key: "tmplUsernameKey", Constraints: ""})
	},
	"NameFirst": func() string {
		return generators["name.first"].Func(Column{Name: "tmplNameFirst", Key: "tmplNameFirstKey", Constraints: ""})
	},
	"NameLast": func() string {
		return generators["name.last"].Func(Column{Name: "tmplNameLast", Key: "tmplNameLastKey", Constraints: ""})
	},
	"Color": func() string {
		return generators["color"].Func(Column{Name: "tmplColor", Key: "tmplColorKey", Constraints: ""})
	},
	"ProductCategory": func() string {
		return generators["product.category"].Func(Column{Name: "tmplProductCategory", Key: "tmplProductCategoryKey", Constraints: ""})
	},
	"ProductName": func() string {
		return generators["product.name"].Func(Column{Name: "tmplProductName", Key: "tmplProductNameKey", Constraints: ""})
	},
	"EventAction": func() string {
		return generators["event.action"].Func(Column{Name: "tmplEventAction", Key: "tmplEventActionKey", Constraints: ""})
	},
	"HTTPMethod": func() string {
		return generators["http.method"].Func(Column{Name: "tmplHTTPMethod", Key: "tmplHTTPMethodKey", Constraints: ""})
	},
	"Name": func() string {
		return generators["name"].Func(Column{Name: "tmplName", Key: "tmplNameKey", Constraints: " "})
	},
	"Email": func() string {
		return generators["email"].Func(Column{Name: "tmplEmail", Key: "tmplEmailKey", Constraints: " "})
	},
	"Domain": func() string {
		return generators["domain"].Func(Column{Name: "tmplDomain", Key: "tmplDomainKey", Constraints: ""})
	},
	"IPv4": func() string {
		return generators["ipv4"].Func(Column{Name: "tmplIPv4", Key: "tmplIPv4Key", Constraints: ""})
	},
	"IPv6": func() string {
		return generators["ipv6"].Func(Column{Name: "tmplIPv6", Key: "tmplIPv6Key", Constraints: ""})
	},
	"MacAddress": func() string {
		return generators["mac.address"].Func(Column{Name: "tmplMacAddress", Key: "tmplMacAddressKey", Constraints: ""})
	},
	"Latitude": func() string {
		return generators["latitude"].Func(Column{Name: "tmplLatitude", Key: "tmplLatitudeKey", Constraints: ""})
	},
	"Longitude": func() string {
		return generators["longitude"].Func(Column{Name: "tmplLongitude", Key: "tmplLongitudeKey", Constraints: ""})
	},
	"Double": func() string {
		return generators["double"].Func(Column{Name: "tmplDouble", Key: "tmplDoubleKey", Constraints: ""})
	},
	"Int": func(params ...int) string {
		constraint := ".."
		a := 0
		b := 1000
		if len(params) == 1 {
			a = params[0]
		}
		if len(params) == 2 {
			a = params[0]
			b = params[1]
		}
		constraint = fmt.Sprintf("%d..%d", a, b)
		return generators["int"].Func(Column{"tmplInt", "tmplIntKey", constraint})
	},
	"Enum": func(keywords ...string) string {
		constraints := strings.Join(keywords, "..")
		return generators["enum"].Func(Column{"tmplEnum", "tmplEnumKey", constraints})
	},
	"File": func(path string) string {
		_, err := os.Stat(path)
		if err != nil {
			fmt.Printf("could not read file: %v", err)
			os.Exit(1)
		}
		return generators["file"].Func(Column{Name: "tmplFile", Key: "tmplFileKey", Constraints: path})
	},
}

func ExecuteTemplate(tmpl string, limit int) (err error) {
	t, err := template.New("template").Funcs(generatorFunctions).Parse(tmpl)
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
