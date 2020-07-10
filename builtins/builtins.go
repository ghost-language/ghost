package builtins

import (
	"fmt"
	"sort"

	"ghostlang.org/ghost/object"
)

var Builtins = map[string]*object.Builtin{
	"first": &object.Builtin{Name: "first", Fn: First},
	"input": &object.Builtin{Name: "input", Fn: Input},
	"last":  &object.Builtin{Name: "last", Fn: Last},
	"len":   &object.Builtin{Name: "len", Fn: Len},
	"print": &object.Builtin{Name: "print", Fn: Print},
	"push":  &object.Builtin{Name: "push", Fn: Push},
	"tail":  &object.Builtin{Name: "tail", Fn: Tail},
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

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
