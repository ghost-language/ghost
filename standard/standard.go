package standard

import (
	"fmt"

	"ghostlang.org/x/ghost/object"
)

var Builtins = map[string]*object.Builtin{}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func RegisterFunction(name string, function object.BuiltinFunction) {
	Builtins[name] = &object.Builtin{Fn: function}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
