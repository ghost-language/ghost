package object

import (
	"github.com/shopspring/decimal"
)

// NUMBER represents the object's type.
const NUMBER = "NUMBER"

// Number objects consist of a decimal value.
type Number struct {
	Value decimal.Decimal
}

// String represents the string form of the number object.
func (n *Number) String() string {
	return n.Value.String()
}

// Type returns the number object type.
func (n *Number) Type() Type {
	return NUMBER
}
