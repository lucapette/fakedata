package fakedata

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

const (
	bashTemplate = `
_fakedata()
{
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="%s"

    if [[ ${cur} == * ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}
complete -F _fakedata fakedata`

	zshTemplate = `
_fakedata () {
    local -a commands
    IFS=$'\n'
    commands=(%s)
    _describe 'arguments' commands
}
compdef _fakedata fakedata`
)

func findCompletionTemplate(sh string) (string, error) {
	switch sh {
	case "bash":
		return bashTemplate, nil

	case "zsh":
		return zshTemplate, nil
	}
	return "", errors.New("Shell could not be found.\nPlease set the $SHELL environment variable and make sure you use one of the supported shells.")
}

func PrintShellCompletionFunction(sh string) (completion string, err error) {
	gens := &bytes.Buffer{}
	allCliArgs := &bytes.Buffer{}

	for _, gen := range NewGenerators() {
		fmt.Fprintf(gens, gen.Name+" ")
	}

	pflag.VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(allCliArgs, "-%s --%s ", f.Shorthand, f.Name)
	})

	t, err := findCompletionTemplate(sh)

	if err != nil {
		return "", err
	}

	cmdList := gens.String() + " " + allCliArgs.String()
	return fmt.Sprintf(t, cmdList), nil
}
