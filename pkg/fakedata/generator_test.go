package fakedata_test

import (
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

var gens = fakedata.NewGenerators()

func BenchmarkNoun(b *testing.B) {
	noun := gens.FindByName("noun")

	for i := 0; i < b.N; i++ {
		noun.Func()
	}
}

func BenchmarkDomain(b *testing.B) {
	domain := gens.FindByName("domain")

	for i := 0; i < b.N; i++ {
		domain.Func()
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

func BenchmarkEmail(b *testing.B) {
	email := gens.FindByName("email")

	for i := 0; i < b.N; i++ {
		email.Func()
	}
}

func BenchmarkName(b *testing.B) {
	name := gens.FindByName("name")

	for i := 0; i < b.N; i++ {
		name.Func()
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
