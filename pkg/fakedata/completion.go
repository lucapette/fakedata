package fakedata

import (
	"bytes"
	"fmt"

	"github.com/spf13/pflag"
)

const bashTemplate = `
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

const zshTemplate = `
_fakedata () {
    local -a commands
    IFS=$'\n'
    commands=(%s)
    _describe 'arguments' commands
}
compdef _fakedata fakedata`

const fishTemplate = `
complete -c fakedata -a '%s'
`

func getTemplate(shell string) (string, error) {
	switch shell {
	case "bash":
		return bashTemplate, nil
	case "zsh":
		return zshTemplate, nil
	case "fish":
		return fishTemplate, nil
	default:
		return "", fmt.Errorf("shell %s not supported. See https://github.com/lucapette/fakedata#completion", shell)
	}
}

// GetCompletionFunc returns a string representing a completion function for the
// given shell. It returns an error for unsupported shell.
func GetCompletionFunc(shell string) (string, error) {
	t, err := getTemplate(shell)
	if err != nil {
		return "", err
	}

	gens := &bytes.Buffer{}
	allCliArgs := &bytes.Buffer{}

	for _, gen := range NewGenerators() {
		fmt.Fprintf(gens, gen.Name+" ") // nolint: errcheck
	}

	pflag.VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(allCliArgs, "--%s ", f.Name) // nolint: errcheck
	})

	cmdList := gens.String() + " " + allCliArgs.String()
	return fmt.Sprintf(t, cmdList), nil
}
