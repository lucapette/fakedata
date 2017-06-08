package fakedata

import (
	"testing"
)

func TestCreateConstraints(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"two integers", []string{"1", "100"}, "1..100"},
		{"three word enum constraints", []string{"feat", "issue", "docs"}, "feat..issue..docs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createConstraints(tt.input); got != tt.want {
				t.Errorf("template.createConstraints = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTemplateNameFromPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"one word template name", "template.tmpl", "template.tmpl"},
		{"absolute template path (unix)", "/var/www/home.tmpl", "home.tmpl"},
		{"absolute template path (Windows)", "C:\\Data\\Templates\\test.tmpl", "test.tmpl"},
		{"relative template path (Unix)", "./views/user/account.tmpl", "account.tmpl"},
		{"relative template path (Windows)", "views\\user\\account.tmpl", "account.tmpl"},
		{"spaced template path (Unix)", "./views/user account/account.tmpl", "account.tmpl"},
		{"spaced template path (Windows)", "views\\user account\\account.tmpl", "account.tmpl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTemplateNameFromPath(tt.input); got != tt.want {
				t.Errorf("template.getTemplateNameFromPath = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseTemplateFromPipe(tt.input, 1); err != nil && (tt.wantErr != true) {
				t.Errorf("template.ParseTemplateFromPipe = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
