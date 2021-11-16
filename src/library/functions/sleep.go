package functions

import (
	"time"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func Sleep(args ...object.Object) object.Object {
	if len(args) != 1 {
		// TODO: error
		return nil
	}

	if args[0].Type() != object.NUMBER {
		// TODO: error
		return nil
	}

	ms := args[0].(*object.Number)
	time.Sleep(time.Duration(ms.Value.IntPart()) * time.Millisecond)

	return value.NULL
}
