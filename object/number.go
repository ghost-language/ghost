package object

import "github.com/shopspring/decimal"

const NUMBER = "NUMBER"

// Number objects consist of a decimal value.
type Number struct {
	Value decimal.Decimal
}

// String represents the number object's value as a string.
func (number *Number) String() string {
	return number.Value.String()
}

// Type returns the number object type.
func (number *Number) Type() Type {
	return NUMBER
}
