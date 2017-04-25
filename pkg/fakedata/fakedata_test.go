package fakedata_test

import (
	"testing"

	"regexp"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

func TestGenerateRow(t *testing.T) {
	type args struct {
		keys   []string
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"email", args{keys: []string{"email"}, format: ""}, `.+?@.+?\..+`},
		{"domain", args{keys: []string{"domain"}, format: ""}, `.+?\..+?`},
		{"username", args{keys: []string{"username"}, format: ""}, `[a-zA-Z0-9]{2,}`},
		{"duoble", args{keys: []string{"double"}, format: ""}, `-?[0-9]+?(\.[0-9]+?)?`},
		{"username domain", args{keys: []string{"username", "domain"}, format: " "}, `[a-zA-Z0-9]{2,} .+?\..+?`},
		{"username domain", args{keys: []string{"username", "domain"}, format: "csv"}, `[a-zA-Z0-9]{2,},.+?\..+?`},
		{"username domain", args{keys: []string{"username", "domain"}, format: "tab"}, `[a-zA-Z0-9]{2,}\t.+?\..+?`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fakedata.GenerateRow(tt.args.keys, tt.args.format)

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
