package object

import (
	"fmt"
	"math"
	"strconv"
)

const NUMBER = "NUMBER"

// Number objects represent numeric values using either int64 or float64 internally.
// Integer operations stay as int64 for speed and exactness.
// Float operations use float64. Division always promotes to float.
type Number struct {
	i       int64
	f       float64
	isFloat bool
}

// NewInt creates a new integer Number.
func NewInt(i int64) *Number {
	return &Number{i: i}
}

// NewFloat creates a new float Number.
func NewFloat(f float64) *Number {
	return &Number{f: f, isFloat: true}
}

// Int64 returns the integer value. If float, truncates.
func (n *Number) Int64() int64 {
	if n.isFloat {
		return int64(n.f)
	}
	return n.i
}

// Float64 returns the float value. If integer, converts.
func (n *Number) Float64() float64 {
	if n.isFloat {
		return n.f
	}
	return float64(n.i)
}

// IsFloat returns true if the number is a float.
func (n *Number) IsFloat() bool {
	return n.isFloat
}

// String represents the number object's value as a string.
func (n *Number) String() string {
	if n.isFloat {
		return strconv.FormatFloat(n.f, 'f', -1, 64)
	}
	return strconv.FormatInt(n.i, 10)
}

// Type returns the number object type.
func (n *Number) Type() Type {
	return NUMBER
}

// MapKey defines a unique hash value for use as a map key.
func (n *Number) MapKey() MapKey {
	if n.isFloat {
		return MapKey{Type: n.Type(), Value: math.Float64bits(n.f)}
	}
	return MapKey{Type: n.Type(), Value: uint64(n.i)}
}

// Method defines the set of methods available on number objects.
func (n *Number) Method(method string, args []Object) (Object, bool) {
	switch method {
	case "round":
		return n.round(args)
	case "floor":
		return n.floor(args)
	case "toString":
		return n.toString(args)
	}

	return nil, false
}

// =============================================================================
// Object methods

func (n *Number) toString(args []Object) (Object, bool) {
	return &String{Value: n.String()}, true
}

func (n *Number) round(args []Object) (Object, bool) {
	places := int64(0)

	if len(args) == 1 {
		if args[0].Type() != NUMBER {
			return nil, false
		}

		places = args[0].(*Number).Int64()
	}

	if !n.isFloat {
		return n, true
	}

	if places == 0 {
		return NewInt(int64(math.Round(n.f))), true
	}

	shift := math.Pow(10, float64(places))
	return NewFloat(math.Round(n.f*shift) / shift), true
}

func (n *Number) floor(args []Object) (Object, bool) {
	if !n.isFloat {
		return n, true
	}

	return NewInt(int64(math.Floor(n.f))), true
}

// =============================================================================
// Arithmetic helpers for the evaluator

// Add returns the sum of two numbers.
func (n *Number) Add(other *Number) *Number {
	if n.isFloat || other.isFloat {
		return NewFloat(n.Float64() + other.Float64())
	}
	return NewInt(n.i + other.i)
}

// Sub returns the difference of two numbers.
func (n *Number) Sub(other *Number) *Number {
	if n.isFloat || other.isFloat {
		return NewFloat(n.Float64() - other.Float64())
	}
	return NewInt(n.i - other.i)
}

// Mul returns the product of two numbers.
func (n *Number) Mul(other *Number) *Number {
	if n.isFloat || other.isFloat {
		return NewFloat(n.Float64() * other.Float64())
	}
	return NewInt(n.i * other.i)
}

// Div always returns a float.
func (n *Number) Div(other *Number) *Number {
	return NewFloat(n.Float64() / other.Float64())
}

// Mod returns the modulo of two numbers.
func (n *Number) Mod(other *Number) *Number {
	if n.isFloat || other.isFloat {
		return NewFloat(math.Mod(n.Float64(), other.Float64()))
	}
	return NewInt(n.i % other.i)
}

// Neg returns the negation.
func (n *Number) Neg() *Number {
	if n.isFloat {
		return NewFloat(-n.f)
	}
	return NewInt(-n.i)
}

// LessThan compares two numbers.
func (n *Number) LessThan(other *Number) bool {
	if n.isFloat || other.isFloat {
		return n.Float64() < other.Float64()
	}
	return n.i < other.i
}

// LessThanOrEqual compares two numbers.
func (n *Number) LessThanOrEqual(other *Number) bool {
	if n.isFloat || other.isFloat {
		return n.Float64() <= other.Float64()
	}
	return n.i <= other.i
}

// GreaterThan compares two numbers.
func (n *Number) GreaterThan(other *Number) bool {
	if n.isFloat || other.isFloat {
		return n.Float64() > other.Float64()
	}
	return n.i > other.i
}

// GreaterThanOrEqual compares two numbers.
func (n *Number) GreaterThanOrEqual(other *Number) bool {
	if n.isFloat || other.isFloat {
		return n.Float64() >= other.Float64()
	}
	return n.i >= other.i
}

// Equal compares two numbers.
func (n *Number) Equal(other *Number) bool {
	if n.isFloat || other.isFloat {
		return n.Float64() == other.Float64()
	}
	return n.i == other.i
}

// Abs returns the absolute value.
func (n *Number) Abs() *Number {
	if n.isFloat {
		return NewFloat(math.Abs(n.f))
	}
	if n.i < 0 {
		return NewInt(-n.i)
	}
	return n
}

// IsNegative returns true if the number is negative.
func (n *Number) IsNeg() bool {
	if n.isFloat {
		return n.f < 0
	}
	return n.i < 0
}

// IsPositive returns true if the number is positive.
func (n *Number) IsPos() bool {
	if n.isFloat {
		return n.f > 0
	}
	return n.i > 0
}

// IsZero returns true if the number is zero.
func (n *Number) IsZero() bool {
	if n.isFloat {
		return n.f == 0
	}
	return n.i == 0
}

// Increment returns n + 1, preserving type.
func (n *Number) Increment() *Number {
	if n.isFloat {
		return NewFloat(n.f + 1)
	}
	return NewInt(n.i + 1)
}

// Decrement returns n - 1, preserving type.
func (n *Number) Decrement() *Number {
	if n.isFloat {
		return NewFloat(n.f - 1)
	}
	return NewInt(n.i - 1)
}

// Cos returns the cosine.
func (n *Number) Cos() *Number {
	return NewFloat(math.Cos(n.Float64()))
}

// Sin returns the sine.
func (n *Number) Sin() *Number {
	return NewFloat(math.Sin(n.Float64()))
}

// Tan returns the tangent.
func (n *Number) Tan() *Number {
	return NewFloat(math.Tan(n.Float64()))
}

// Stringer for debugging.
func (n *Number) GoString() string {
	if n.isFloat {
		return fmt.Sprintf("Number(float:%g)", n.f)
	}
	return fmt.Sprintf("Number(int:%d)", n.i)
}
