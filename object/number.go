package object

import (
	"fmt"
	"math"
	"strconv"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

// Number objects represent numeric values using either int64 or float64 internally.
// Integer operations stay as int64 for speed and exactness.
// Float operations use float64. Division always promotes to float.
type Number struct {
	i       int64
	f       float64
	isFloat bool
}

// Bounds of the preallocated small-integer cache. Loop counters, indices, and
// list lengths overwhelmingly fall in this range, so interning them removes a
// heap allocation from nearly every arithmetic operation. Number is immutable
// (its fields are unexported and never reassigned after construction), which is
// what makes sharing these instances safe.
const (
	smallIntMin = -128
	smallIntMax = 1024
)

var smallInts [smallIntMax - smallIntMin + 1]Number

func init() {
	for index := range smallInts {
		smallInts[index] = Number{i: int64(index + smallIntMin)}
	}
}

// NewInt creates a new integer Number, returning a shared instance for small
// values.
func NewInt(i int64) *Number {
	if i >= smallIntMin && i <= smallIntMax {
		return &smallInts[i-smallIntMin]
	}

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
func (n *Number) Method(method string, tok token.Token, args []Object) (Object, bool) {
	switch method {
	case "round":
		return n.round(tok, args)
	case "floor":
		return n.floor(tok, args)
	case "ceil":
		return n.ceil(tok, args)
	case "abs":
		return n.abs(tok, args)
	case "pow":
		return n.pow(tok, args)
	case "clamp":
		return n.clamp(tok, args)
	case "isNaN":
		return n.isNaNMethod(tok, args)
	case "isFinite":
		return n.isFiniteMethod(tok, args)
	case "isInfinite":
		return n.isInfiniteMethod(tok, args)
	case "isInteger":
		return n.isIntegerMethod(tok, args)
	case "isEven":
		return n.isEvenMethod(tok, args)
	case "isOdd":
		return n.isOddMethod(tok, args)
	case "isNegative":
		return n.isNegativeMethod(tok, args)
	case "isPositive":
		return n.isPositiveMethod(tok, args)
	case "isZero":
		return n.isZeroMethod(tok, args)
	case "toString":
		return n.toString(tok, args)
	}

	return nil, false
}

// =============================================================================
// Object methods

func (n *Number) toString(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.toString()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: n.String()}, true
}

func (n *Number) round(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("number.round()", tok, args, 0, 1); err != nil {
		return err, true
	}

	places := int64(0)

	if len(args) == 1 {
		digits, err := NumberArgument("number.round()", tok, args, 0)

		if err != nil {
			return err, true
		}

		places = digits.Int64()
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

func (n *Number) floor(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.floor()", tok, args, 0); err != nil {
		return err, true
	}

	if !n.isFloat {
		return n, true
	}

	return NewInt(int64(math.Floor(n.f))), true
}

func (n *Number) ceil(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.ceil()", tok, args, 0); err != nil {
		return err, true
	}

	if !n.isFloat {
		return n, true
	}

	return NewInt(int64(math.Ceil(n.f))), true
}

func (n *Number) abs(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.abs()", tok, args, 0); err != nil {
		return err, true
	}

	return n.Abs(), true
}

// pow keeps whole numbers whole where it can, mirroring math.pow: an integer
// base raised to a non-negative integer exponent has an exact integer answer,
// so `4.pow(2)` stays a Number that can index a list. When the answer would
// not fit an int64, or either operand is a float, it falls back to floating
// point rather than wrapping.
func (n *Number) pow(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.pow()", tok, args, 1); err != nil {
		return err, true
	}

	exponent, err := NumberArgument("number.pow()", tok, args, 0)

	if err != nil {
		return err, true
	}

	if !n.isFloat && !exponent.isFloat && !exponent.IsNeg() {
		if result, ok := integerPower(n.i, exponent.i); ok {
			return NewInt(result), true
		}
	}

	return NewFloat(math.Pow(n.Float64(), exponent.Float64())), true
}

// clamp keeps the receiver inside [low, high], answering with one of the
// three values it was given rather than a computed one, so clamping whole
// numbers leaves them whole.
func (n *Number) clamp(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.clamp()", tok, args, 2); err != nil {
		return err, true
	}

	low, err := NumberArgument("number.clamp()", tok, args, 0)

	if err != nil {
		return err, true
	}

	high, err := NumberArgument("number.clamp()", tok, args, 1)

	if err != nil {
		return err, true
	}

	if low.GreaterThan(high) {
		return NewError(fault.Value, tok, "`number.clamp()` expects the lower bound to be no greater than the upper bound"), true
	}

	if n.LessThan(low) {
		return low, true
	}

	if n.GreaterThan(high) {
		return high, true
	}

	return n, true
}

func (n *Number) isNaNMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isNaN()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.isFloat && math.IsNaN(n.f)}, true
}

func (n *Number) isInfiniteMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isInfinite()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.isFloat && math.IsInf(n.f, 0)}, true
}

func (n *Number) isFiniteMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isFinite()", tok, args, 0); err != nil {
		return err, true
	}

	finite := !(n.isFloat && math.IsNaN(n.f)) && !(n.isFloat && math.IsInf(n.f, 0))

	return &Boolean{Value: finite}, true
}

func (n *Number) isIntegerMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isInteger()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.isIntegerValue()}, true
}

func (n *Number) isEvenMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isEven()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.isIntegerValue() && math.Mod(n.Float64(), 2) == 0}, true
}

func (n *Number) isOddMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isOdd()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.isIntegerValue() && math.Abs(math.Mod(n.Float64(), 2)) == 1}, true
}

func (n *Number) isNegativeMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isNegative()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.IsNeg()}, true
}

func (n *Number) isPositiveMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isPositive()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.IsPos()}, true
}

func (n *Number) isZeroMethod(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("number.isZero()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: n.IsZero()}, true
}

// isIntegerValue reports whether n holds a mathematically whole value: every
// non-float Number does, and a float only if it is finite and has no
// fractional part.
func (n *Number) isIntegerValue() bool {
	if !n.isFloat {
		return true
	}

	return !math.IsNaN(n.f) && !math.IsInf(n.f, 0) && n.f == math.Trunc(n.f)
}

// integerPower computes base^exponent for a non-negative exponent, reporting
// false if the exact int64 result would overflow.
func integerPower(base int64, exponent int64) (int64, bool) {
	result := int64(1)

	for exponent > 0 {
		if exponent&1 == 1 {
			product, ok := multiplyChecked(result, base)

			if !ok {
				return 0, false
			}

			result = product
		}

		exponent >>= 1

		if exponent == 0 {
			break
		}

		square, ok := multiplyChecked(base, base)

		if !ok {
			return 0, false
		}

		base = square
	}

	return result, true
}

// multiplyChecked multiplies two int64s, reporting false if the exact result
// would overflow rather than silently wrapping.
func multiplyChecked(left int64, right int64) (int64, bool) {
	if left == 0 || right == 0 {
		return 0, true
	}

	if (left == -1 && right == math.MinInt64) || (right == -1 && left == math.MinInt64) {
		return 0, false
	}

	product := left * right

	if product/right != left {
		return 0, false
	}

	return product, true
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
