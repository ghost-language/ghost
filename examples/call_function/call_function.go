package main

import (
	"fmt"

	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/object"
)

func main() {
	ghost.RegisterFunction("write", writeFunction)

	ghost.NewScript(`write("this is a custom function."); foobar := "example string"`)

	env := ghost.Evaluate()

	ghost.Call(`write("this was called separately."); write(foobar); foobar := "crash override"`, env)
	ghost.Call(`write(foobar)`, env)
}

func writeFunction(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return ghost.NewError("wrong number of arguments. got=%d, expected=1", len(args))
	}

	fmt.Println(args[0].Inspect())

	return nil
}
