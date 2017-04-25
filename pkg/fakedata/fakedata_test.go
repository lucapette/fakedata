package fakedata_test

import (
	"testing"

	"regexp"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

func TestGenerateRow(t *testing.T) {
	csv := fakedata.NewSeparatorFormatter(",")
	def := fakedata.NewSeparatorFormatter(" ")
	tab := fakedata.NewSeparatorFormatter("\t")

	type args struct {
		columns   fakedata.Columns
		formatter fakedata.Formatter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"email", args{columns: fakedata.Columns{{Key: "email"}}, formatter: def}, `.+?@.+?\..+`},
		{"domain", args{columns: fakedata.Columns{{Key: "domain"}}, formatter: def}, `.+?\..+?`},
		{"username", args{columns: fakedata.Columns{{Key: "username"}}, formatter: def}, `[a-zA-Z0-9]{2,}`},
		{"duoble", args{columns: fakedata.Columns{{Key: "double"}}, formatter: def}, `-?[0-9]+?(\.[0-9]+?)?`},
		{"username domain", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, formatter: def}, `[a-zA-Z0-9]{2,} .+?\..+?`},
		{"username domain csv", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, formatter: csv}, `[a-zA-Z0-9]{2,},.+?\..+?`},
		{"username domain tab", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, formatter: tab}, `[a-zA-Z0-9]{2,}\t.+?\..+?`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fakedata.GenerateRow(tt.args.columns, tt.args.formatter)

			matched, err := regexp.MatchString(tt.want, got)
			if err != nil {
				t.Error(err.Error())
			}

			if !matched {
				t.Errorf("GenerateRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
