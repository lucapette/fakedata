package fakedata

import (
	"testing"
	"text/template"
)

func TestParseTemplateFromPipe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"name generator", "{{ Name }}", false},
		{"undefined function", "{{ NotAGen }}", true},
		{"invalid template function, unexpected end", "{{ Name }} {{ end }}", true},
		{"invalid template function, unclosed range", "{{ range Loop 10 }} {{ Name }}, {{ Int }}", true},
		{"printf in template", "{{ printf \"%s, %s\" NameLast NameFirst }}", false},
		{"function with arguments, int", "{{ Int 12 15 }}", false},
		{"function with arguments, enum", "{{ Enum \"Feature\" \"Issue\" \"Resolved\" \"On hold\" }}", false},
		{"no file", "{{ File \"NotAFile.txt\" }}", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseTemplateFromPipe(tt.input); err != nil && (tt.wantErr != true) {
				t.Errorf("template.ParseTemplateFromPipe = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteTemplate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple template", "{{ Name }}", false},
		{"two generators", "{{ Name }} {{ Date }}", false},
		{"undefined function", "{{ GetDataFromUrl http://fake.link/ }}", true},
		{"printf in template", "{{ printf \"%s %s\" NameLast NameFirst }}", false},
		{"unclosed range", "{{range Loop 5}} {{Name}}", true},
		{"template variables", "{{$a := Name }} {{ $a }}", false},
		{"function with parameters", "{{ Enum \"Lorem\" \"Ipsum\"}}", false},
		{"invalid template", "{{ Name, {{ Int 15 20 }},  }}", true},
		{"text template", "Hello, World!", false},
		{"no file", "{{ File \"NotAFile.txt\" }}", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := template.New(tt.name).Funcs(generatorFunctions).Parse(tt.input)
			if err != nil && !tt.wantErr {
				t.Errorf("template.ExecuteTemplate error = %v, want %v", err, tt.wantErr)
			}

		})
	}
}
