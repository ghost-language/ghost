package builtins

import (
	"os"

	"ghostlang.org/x/ghost/object"
)

func init() {
	RegisterFunction("os.args", osArgsFunction)
}

// osArgsFunction returns the arguments passed to the program
func osArgsFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	list := &object.List{}
	arguments := os.Args[2:]

	for _, argument := range arguments {
		list.Elements = append(list.Elements, &object.String{Value: argument})
	}

	return list
}
