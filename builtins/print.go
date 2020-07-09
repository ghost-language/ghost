package builtins

import (
	"fmt"

	"ghostlang.org/ghost/object"
)

func Print(args ...object.Object) object.Object {
	if len(args) > 0 {
		fmt.Println(args[0].Inspect())
	} else {
		fmt.Println()
	}

	return nil
}
