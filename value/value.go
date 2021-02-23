package value

import "ghostlang.org/x/ghost/object"

var (
	// NULL represents a null value.
	NULL = &object.Null{}

	// TRUE represent a true value.
	TRUE = &object.Boolean{Value: true}

	// FALSE represents a false value.
	FALSE = &object.Boolean{Value: false}
)
