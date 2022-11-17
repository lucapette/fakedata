package fakedata_test

import (
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

var gens = fakedata.NewGenerators()

func BenchmarkGenerators(b *testing.B) {
	for i := 0; i < len(gens); i++ {
		g := gens[i]

		if !g.IsCustom() {
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