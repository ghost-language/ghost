package ast

import "ghostlang.org/x/ghost/token"

type Map struct {
	ExpressionNode
	Token token.Token

	// Pairs is ordered the way it was written, not keyed by the expression
	// itself the way a Go map would be - a bare map here would forget source
	// order before evaluation ever saw it, undermining Map's own
	// insertion-order guarantee (§13.5, §14 decision 2) at its earliest
	// possible source. Two different key expressions that happen to
	// evaluate to the same map key (`{x: 1, x: 2}`) are both kept; which one
	// wins is decided at evaluation, the same as any other repeated
	// assignment - last one written.
	Pairs []MapEntry
}

// MapEntry is one `key: value` pair of a map literal, in the order it was
// written.
type MapEntry struct {
	Key   ExpressionNode
	Value ExpressionNode
}
