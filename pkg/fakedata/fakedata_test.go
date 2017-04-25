package fakedata_test

import (
	"testing"

	"regexp"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

func TestGenerateRow(t *testing.T) {
	type args struct {
		columns fakedata.Columns
		format  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"email", args{columns: fakedata.Columns{{Key: "email"}}, format: ""}, `.+?@.+?\..+`},
		{"domain", args{columns: fakedata.Columns{{Key: "domain"}}, format: ""}, `.+?\..+?`},
		{"username", args{columns: fakedata.Columns{{Key: "username"}}, format: ""}, `[a-zA-Z0-9]{2,}`},
		{"duoble", args{columns: fakedata.Columns{{Key: "double"}}, format: ""}, `-?[0-9]+?(\.[0-9]+?)?`},
		{"username domain", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, format: " "}, `[a-zA-Z0-9]{2,} .+?\..+?`},
		{"username domain", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, format: "csv"}, `[a-zA-Z0-9]{2,},.+?\..+?`},
		{"username domain", args{columns: fakedata.Columns{{Key: "username"}, {Key: "domain"}}, format: "tab"}, `[a-zA-Z0-9]{2,}\t.+?\..+?`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fakedata.GenerateRow(tt.args.columns, tt.args.format)

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
