package ast

import "ghostlang.org/x/ghost/token"

type Import struct {
	ExpressionNode
	Token token.Token
	Path  *String
	Alias *Identifier

	// Identifiers and Everything are set only by the JS-style combined form,
	// `import name, { a, b } from "path"` (or `import name, { * } from
	// "path"`) — named exports pulled from the same module alongside the
	// whole-module binding above (Alias carries `name` there, exactly as an
	// `as` alias would). A bare `import "path"` leaves both zero.
	Identifiers map[string]*Identifier
	Everything  bool
}
