package value

import (
	"ghostlang.org/x/ghost/object"
)

var (
	// TRUE represents a true value.
	TRUE = &object.Boolean{Value: true}

	// FALSE represents a false value.
	FALSE = &object.Boolean{Value: false}

	// NULL represents a null value.
	NULL = &object.Null{}

	// BREAK represents a break value.
	BREAK = &object.Break{}

	// CONTINUE represents a continue value.
	CONTINUE = &object.Continue{}
)
