package fakedata_test

import (
	"testing"

	"github.com/lucapette/fakedata/pkg/fakedata"
)

func TestPrintShellCompletionFunction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"zsh shell", "zsh", false},
		{"bash shell", "bash", false},
		{"fish shell", "fish", true},
		{"bsh shell (spelling mistake)", "bsh", true},
		{"empty", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := fakedata.GetCompletionFunc(tc.input); err != nil && !tc.wantErr {
				t.Errorf("expected err to be %v but got %v", tc.wantErr, err)
			}
		})
	}
}
