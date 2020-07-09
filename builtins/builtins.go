package builtins

import (
	"fmt"
	"sort"

	"ghostlang.org/ghost/object"
)

var Builtins = map[string]*object.Builtin{
	"input": &object.Builtin{Name: "input", Fn: Input},
	"len":   &object.Builtin{Name: "len", Fn: Len},
	"print": &object.Builtin{Name: "print", Fn: Print},
}

var BuiltinsIndex []*object.Builtin

func init() {
	var keys []string
	for k := range Builtins {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		BuiltinsIndex = append(BuiltinsIndex, Builtins[k])
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
