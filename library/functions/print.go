package functions

import (
	"fmt"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func Print(args ...object.Object) object.Object {
	fmt.Println("print()")

	return value.TRUE
}
