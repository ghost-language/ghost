package object

import (
	"strconv"

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

func (n *Number) Set(obj Object) {
	n.Value = obj.(*Number).Value
}

func (n *Number) MapKey() MapKey {
	value, _ := strconv.ParseUint(n.Value.String(), 10, 64)

	return MapKey{Type: n.Type(), Value: value}
}
