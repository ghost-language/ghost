package ast

import "ghostlang.org/x/ghost/token"

// ListPattern is the assignment target of a list-destructuring assignment,
// `[a, b] = list`: positional, so Targets[i] binds to the value at index i
// of whatever the right-hand side evaluates to.
type ListPattern struct {
	AssignmentNode
	Token   token.Token
	Targets []*Identifier
}

// MapPattern is the assignment target of a map-destructuring assignment,
// `{x, y} = map` (shorthand, matching map literal shorthand keys) or
// `{x: a} = map` (binds map key `x` to the local name `a`). Each pair reads
// Source from the map on the right and binds it to Target.
type MapPattern struct {
	AssignmentNode
	Token token.Token
	Pairs []MapPatternPair
}

type MapPatternPair struct {
	Source *Identifier
	Target *Identifier
}
