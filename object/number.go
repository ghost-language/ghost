package object

import (
	"ghostlang.org/ghost/decimal"
)

type Number struct {
	Value decimal.Decimal
}

func (n *Number) Inspect() string {
	return n.Value.String()
}

func (n *Number) Type() ObjectType {
	return NUMBER_OBJ
}
