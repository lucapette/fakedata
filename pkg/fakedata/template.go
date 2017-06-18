package fakedata

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// generatorFunctions holds all the functions available for the template
var generatorFunctions = template.FuncMap{
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
		return generators["int"].Func(Column{"Int", "int", constraint})
	},
	"Enum": func(keywords ...string) string {
		constraints := createConstraints(keywords)
		return generators["enum"].Func(Column{"Enum", "enu", constraints})
	},
	"File": func(path string) string {
		return generators["file"].Func(Column{Name: "tmplA", Key: "tmplAKey", Constraints: path})
	},
}

func createConstraints(params []string) string {
	return strings.Join(params, "..")
}

func getTemplateNameFromPath(name string) string {
	ts := strings.FieldsFunc(name, splitPathName)
	tn := ts[len(ts)-1]
	return tn
}

// this custom split function is used with strings.FieldsFunc to split the path
// by `/` (Unix, MacOS) or `\` (Windows) for absolute and relative paths to template files
func splitPathName(r rune) bool {
	return r == '/' || r == '\\'
}

// ParseTemplate takes a path to a template file as argument. It parses the template file and executes it on
// os.Stdout.
func ParseTemplate(path string) (tmp *template.Template, err error) {
	tn := getTemplateNameFromPath(path)
	tmp, err = template.New(tn).Funcs(generatorFunctions).ParseFiles(path)
	if err != nil {
		return nil, err
	}

	return tmp, nil
}

// ParseTemplateFromPipe takes a string as template, parses it and executed the template. The function returns an error
// or nil on success. The template is written to os.Stdout
func ParseTemplateFromPipe(t string) (tmp *template.Template, err error) {
	tmp, err = template.New("stdin").Funcs(generatorFunctions).Parse(t)

	if err != nil {
		return nil, err
	}

	return tmp, nil
}

func ExecuteTemplate(t *template.Template, limit int) (err error) {
	b := io.Writer(os.Stdout)
	for i := 1; i <= limit; i++ {
		err = t.Execute(b, nil)
	}
	return err
}
