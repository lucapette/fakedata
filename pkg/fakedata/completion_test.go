package fakedata

import (
	"testing"
)

func TestPrintShellCompletionFunction(t *testing.T) {
	var tests = []struct {
		Name    string
		Input   string
		WantErr bool
	}{
		{"zsh shell", "zsh", false},
		{"bash shell", "bash", false},
		{"fish shell", "fish", true},
		{"bsh shell (spelling mistake)", "bsh", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := PrintShellCompletionFunction(tt.Input)
			if err != nil && !tt.WantErr {
				t.Errorf("Shell Completion error. Got %v but want %v", err, tt.WantErr)
			}
		})
	}
}
