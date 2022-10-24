package testutil

import "github.com/kr/pretty"

func Diff(expected, actual interface{}) []string {
	return pretty.Diff(expected, actual)
}
