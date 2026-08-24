package ast

import "ghostlang.org/x/ghost/token"

type Import struct {
	ExpressionNode
	Token token.Token
	Path  *String
	Alias *Identifier

	// Identifiers and Everything are set only by the combined form,
	// `import "path", { a, b }` (or `import "path", { * }`) — named exports
	// pulled from the same module alongside the whole-module binding above.
	// A bare `import "path"` leaves both zero.
	Identifiers map[string]*Identifier
	Everything  bool
}
