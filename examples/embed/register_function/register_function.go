package main

import (
	"fmt"

	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/object"
)

func main() {
	ghost.RegisterFunction("write", writeFunction)

	ghost.NewScript(`write("this is a custom function.")`)

	ghost.Evaluate()
}

func writeFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return ghost.NewError(line, "wrong number of arguments. got=%d, expected=1", len(args))
	}

	fmt.Println(args[0].Inspect())

	return nil
}
