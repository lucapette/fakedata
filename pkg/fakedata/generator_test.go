package fakedata_test

import (
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

var gens = fakedata.NewGenerators()

func BenchmarkGenerators(b *testing.B) {
	for i := 0; i < len(gens); i++ {
		g := gens[i]

		if !g.IsCustom() && !g.Hidden {
			b.Run(g.Name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					g.Func()
				}
			})
		}
	}
}

func BenchmarkEnum(b *testing.B) {
	enum := gens.FindByName("enum")

	enumFunc, err := enum.CustomFunc("")
	if err != nil {
		b.Fatalf("cannot create enum: %s", err)
	}

	for i := 0; i < b.N; i++ {
		enumFunc()
	}
}

func BenchmarkInt(b *testing.B) {
	integer := gens.FindByName("int")

	integerFunc, err := integer.CustomFunc("10000000,9999999999")
	if err != nil {
		b.Fatalf("cannot create int: %s", err)
	}

	for i := 0; i < b.N; i++ {
		integerFunc()
	}
}

func BenchmarkPhoneLocal(b *testing.B) {
	phoneLocal := gens.FindByName("phone.local")

	phoneLocalFunc, err := phoneLocal.CustomFunc("11")
	if err != nil {
		b.Fatalf("cannot create phone.local: %s", err)
	}

	for i := 0; i < b.N; i++ {
		phoneLocalFunc()
	}
}

func BenchmarkFile(b *testing.B) {
	file := gens.FindByName("file")

	fileFunc, err := file.CustomFunc("../../testutil/fixtures/file.txt")
	if err != nil {
		b.Fatalf("cannot open fixture: %s", err)
	}

	for i := 0; i < b.N; i++ {
		fileFunc()
	}
}
