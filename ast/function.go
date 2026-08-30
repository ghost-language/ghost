package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Function struct {
	ExpressionNode
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Defaults   map[string]ExpressionNode
	Body       *Block

	// Rest marks the last entry of Parameters as a rest parameter (`...args`):
	// it collects every argument from its position onward into a list, rather
	// than binding a single value. A rest parameter never has a default - it
	// is always at least an empty list.
	Rest bool
}
