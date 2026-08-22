package modules

import (
	"math"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

// The math module covers three layers that build on one another. At the bottom
// are the scalar operations any standard library is expected to carry: roots,
// logarithms, trigonometry, rounding, predicates. Above those sits broadcasting,
// which lifts every one of them to work on lists, and on lists of lists, without
// the operation itself knowing that lists exist. On top of that sit the array
// methods - statistics, reductions, linear algebra - which read whole
// collections at once. Because the middle layer is general, `math.sqrt` means
// the same thing whether it is handed a number, a vector, or a matrix.
//
// The array layers live in math_statistics.go and math_array.go.

var MathMethods = map[string]*object.LibraryFunction{}
var MathProperties = map[string]*object.LibraryProperty{}

// defaultTolerance is how close two numbers must be for math.isClose to call
// them equal. It is loose enough to absorb the rounding that accumulates over a
// handful of floating point operations, and tight enough that genuinely
// different values are still reported as different.
const defaultTolerance = 1e-9

func init() {
	// Sign and rounding.
	registerElementwise("abs", 1, unaryNumber((*object.Number).Abs))
	registerElementwise("sign", 1, unaryNumber(numberSign))
	registerElementwise("floor", 1, unaryInteger(math.Floor))
	registerElementwise("ceil", 1, unaryInteger(math.Ceil))
	registerElementwise("truncate", 1, unaryInteger(math.Trunc))
	RegisterMethod(MathMethods, "round", mathRound)

	// Powers, roots, and logarithms.
	registerElementwise("sqrt", 1, unaryFloat(math.Sqrt))
	registerElementwise("cbrt", 1, unaryFloat(math.Cbrt))
	registerElementwise("square", 1, unaryNumber(numberSquare))
	registerElementwise("reciprocal", 1, numberReciprocal)
	registerElementwise("exp", 1, unaryFloat(math.Exp))
	registerElementwise("exp2", 1, unaryFloat(math.Exp2))
	registerElementwise("expm1", 1, unaryFloat(math.Expm1))
	registerElementwise("log2", 1, unaryFloat(math.Log2))
	registerElementwise("log10", 1, unaryFloat(math.Log10))
	registerElementwise("log1p", 1, unaryFloat(math.Log1p))
	registerElementwise("pow", 2, numberPower)
	registerElementwise("hypot", 2, binaryFloat(math.Hypot))
	RegisterMethod(MathMethods, "log", mathLog)

	// Trigonometry. Angles are radians, as they are everywhere in Go's math
	// package; degrees and radians convert between the two.
	registerElementwise("sin", 1, unaryFloat(math.Sin))
	registerElementwise("cos", 1, unaryFloat(math.Cos))
	registerElementwise("tan", 1, unaryFloat(math.Tan))
	registerElementwise("asin", 1, unaryFloat(math.Asin))
	registerElementwise("acos", 1, unaryFloat(math.Acos))
	registerElementwise("atan", 1, unaryFloat(math.Atan))
	registerElementwise("atan2", 2, binaryFloat(math.Atan2))
	registerElementwise("sinh", 1, unaryFloat(math.Sinh))
	registerElementwise("cosh", 1, unaryFloat(math.Cosh))
	registerElementwise("tanh", 1, unaryFloat(math.Tanh))
	registerElementwise("asinh", 1, unaryFloat(math.Asinh))
	registerElementwise("acosh", 1, unaryFloat(math.Acosh))
	registerElementwise("atanh", 1, unaryFloat(math.Atanh))
	registerElementwise("degrees", 1, unaryFloat(toDegrees))
	registerElementwise("radians", 1, unaryFloat(toRadians))

	// Special functions, for the statistics and error terms that come up often
	// enough to be worth having and are miserable to write by hand.
	registerElementwise("gamma", 1, unaryFloat(math.Gamma))
	registerElementwise("logGamma", 1, unaryFloat(logGamma))
	registerElementwise("erf", 1, unaryFloat(math.Erf))
	registerElementwise("erfc", 1, unaryFloat(math.Erfc))

	// Arithmetic as methods. The operators already add two numbers; these add a
	// number to every element of a list, or two lists to each other, which is
	// what makes vector and matrix work read as arithmetic.
	registerElementwise("add", 2, binaryNumber((*object.Number).Add))
	registerElementwise("subtract", 2, binaryNumber((*object.Number).Sub))
	registerElementwise("multiply", 2, binaryNumber((*object.Number).Mul))
	registerElementwise("divide", 2, numberDivide)
	registerElementwise("mod", 2, numberModulo)
	registerElementwise("remainder", 2, binaryFloat(math.Remainder))
	registerElementwise("copySign", 2, binaryFloat(math.Copysign))
	registerElementwise("maximum", 2, binaryNumber(numberMaximum))
	registerElementwise("minimum", 2, binaryNumber(numberMinimum))

	// Predicates.
	registerElementwise("isNaN", 1, unaryPredicate(isNotANumber))
	registerElementwise("isFinite", 1, unaryPredicate(isFinite))
	registerElementwise("isInfinite", 1, unaryPredicate(isInfinite))
	registerElementwise("isInteger", 1, unaryPredicate(isInteger))
	registerElementwise("isEven", 1, unaryPredicate(isEven))
	registerElementwise("isOdd", 1, unaryPredicate(isOdd))
	registerElementwise("isNegative", 1, unaryPredicate((*object.Number).IsNeg))
	registerElementwise("isPositive", 1, unaryPredicate((*object.Number).IsPos))
	registerElementwise("isZero", 1, unaryPredicate((*object.Number).IsZero))
	RegisterMethod(MathMethods, "isClose", mathIsClose)

	// Bounds and interpolation.
	registerElementwise("clamp", 3, numberClamp)
	registerElementwise("lerp", 3, ternaryFloat(interpolate))
	registerElementwise("smoothstep", 3, ternaryFloat(smoothstep))
	RegisterMethod(MathMethods, "noise", mathNoise)

	// Random numbers. These share the generator behind the random module, so a
	// single seed governs both and a seeded run stays reproducible no matter
	// which of the two a program reaches for.
	RegisterMethod(MathMethods, "random", mathRandom)
	RegisterMethod(MathMethods, "randomSeed", mathRandomSeed)

	registerMathStatistics()
	registerMathArrays()

	// Constants.
	registerConstant("pi", math.Pi)
	registerConstant("tau", 2*math.Pi)
	registerConstant("e", math.E)
	registerConstant("phi", math.Phi)
	registerConstant("sqrt2", math.Sqrt2)
	registerConstant("sqrtPi", math.SqrtPi)
	registerConstant("ln2", math.Ln2)
	registerConstant("ln10", math.Ln10)
	registerConstant("log2e", math.Log2E)
	registerConstant("log10e", math.Log10E)
	registerConstant("epsilon", math.Nextafter(1, 2)-1)
	registerConstant("smallestNumber", math.SmallestNonzeroFloat64)
	registerConstant("largestNumber", math.MaxFloat64)
	registerConstant("infinity", math.Inf(1))
	registerConstant("nan", math.NaN())
	registerIntegerConstant("largestInteger", math.MaxInt64)
	registerIntegerConstant("smallestInteger", math.MinInt64)
}

// =============================================================================
// Broadcasting

// elementwiseOp is a math operation written against plain numbers. Registering
// one with registerElementwise is what lifts it to lists: the operation is
// written once for numbers and works on vectors and matrices for free.
type elementwiseOp func(tok token.Token, values []*object.Number) object.Object

// broadcast applies an operation across its arguments elementwise, stretching
// the smaller shapes across the larger. The rules live in object.Broadcast,
// which the evaluator also uses for `+` and `*` on lists, so `math.add(a, b)`
// and `a + b` are one operation reached two ways rather than two that have to
// be kept in step.
func broadcast(name string, tok token.Token, args []object.Object, operation elementwiseOp) object.Object {
	result, mismatch := object.Broadcast(args, func(values []*object.Number) object.Object {
		return operation(tok, values)
	})

	if mismatch != nil {
		return object.NewError(fault.Value, tok, "`%s()` %s", name, mismatch.Reason)
	}

	return result
}

// elementwise turns an operation into a library method of fixed arity.
func elementwise(name string, count int, operation elementwiseOp) object.GoFunction {
	return func(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
		if err := arity(name, tok, args, count); err != nil {
			return err
		}

		return broadcast(name, tok, args, operation)
	}
}

func registerElementwise(name string, count int, operation elementwiseOp) {
	RegisterMethod(MathMethods, name, elementwise("math."+name, count, operation))
}

func registerConstant(name string, constant float64) {
	RegisterProperty(MathProperties, name, func(scope *object.Scope, tok token.Token) object.Object {
		return object.NewFloat(constant)
	})
}

func registerIntegerConstant(name string, constant int64) {
	RegisterProperty(MathProperties, name, func(scope *object.Scope, tok token.Token) object.Object {
		return object.NewInt(constant)
	})
}

// =============================================================================
// Adapters, which spell out how a Go operation becomes an elementwiseOp.

// unaryNumber keeps a whole number whole, for the operations whose result is
// exact when their input is.
func unaryNumber(operation func(*object.Number) *object.Number) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return operation(values[0])
	}
}

func unaryFloat(operation func(float64) float64) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return object.NewFloat(operation(values[0].Float64()))
	}
}

// unaryInteger is for the rounding operations, whose results should come back
// as whole numbers so they can index a list without a further conversion.
func unaryInteger(operation func(float64) float64) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return object.NewInt(int64(operation(values[0].Float64())))
	}
}

func unaryPredicate(operation func(*object.Number) bool) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return toBoolean(operation(values[0]))
	}
}

func binaryNumber(operation func(*object.Number, *object.Number) *object.Number) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return operation(values[0], values[1])
	}
}

func binaryFloat(operation func(float64, float64) float64) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return object.NewFloat(operation(values[0].Float64(), values[1].Float64()))
	}
}

func ternaryFloat(operation func(float64, float64, float64) float64) elementwiseOp {
	return func(tok token.Token, values []*object.Number) object.Object {
		return object.NewFloat(operation(values[0].Float64(), values[1].Float64(), values[2].Float64()))
	}
}

// =============================================================================
// Methods whose argument count varies, and so cannot be a table entry.

// mathRound rounds to the nearest whole number, or to a number of decimal
// places when given a second argument. Rounding to zero places answers with a
// whole number; rounding to more keeps the fraction it was asked to keep.
func mathRound(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.round", tok, args, 1, 2); err != nil {
		return err
	}

	places := int64(0)

	if len(args) == 2 {
		given, err := integerAt("math.round", tok, args, 1)

		if err != nil {
			return err
		}

		places = given
	}

	if places == 0 {
		return broadcast("math.round", tok, args[:1], unaryInteger(math.Round))
	}

	shift := math.Pow(10, float64(places))

	return broadcast("math.round", tok, args[:1], unaryFloat(func(given float64) float64 {
		return math.Round(given*shift) / shift
	}))
}

// mathLog returns the natural logarithm, or the logarithm in the given base.
func mathLog(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.log", tok, args, 1, 2); err != nil {
		return err
	}

	if len(args) == 1 {
		return broadcast("math.log", tok, args, unaryFloat(math.Log))
	}

	base, err := floatAt("math.log", tok, args, 1)

	if err != nil {
		return err
	}

	return broadcast("math.log", tok, args[:1], unaryFloat(func(given float64) float64 {
		return math.Log(given) / math.Log(base)
	}))
}

// mathIsClose reports whether two numbers are equal to within a tolerance. It
// is the comparison to reach for after any run of floating point arithmetic,
// where `==` asks a question the arithmetic cannot answer.
func mathIsClose(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.isClose", tok, args, 2, 3); err != nil {
		return err
	}

	operands := make([]object.Object, 3)
	operands[0] = args[0]
	operands[1] = args[1]
	operands[2] = object.NewFloat(defaultTolerance)

	if len(args) == 3 {
		if _, err := numberAt("math.isClose", tok, args, 2); err != nil {
			return err
		}

		operands[2] = args[2]
	}

	return broadcast("math.isClose", tok, operands, func(tok token.Token, values []*object.Number) object.Object {
		return toBoolean(closeEnough(values[0].Float64(), values[1].Float64(), values[2].Float64()))
	})
}

// mathNoise samples value noise smoothed with a cubic curve: continuous,
// repeatable for the same input, and in the 0-1 range. Where random numbers
// jitter, noise drifts, which is what makes it the right source for terrain
// heights, organic-looking variation, and anything that should wander rather
// than jump.
func mathNoise(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.noise", tok, args, 1, 2); err != nil {
		return err
	}

	x, err := floatAt("math.noise", tok, args, 0)

	if err != nil {
		return err
	}

	y := 0.0

	if len(args) == 2 {
		given, err := floatAt("math.noise", tok, args, 1)

		if err != nil {
			return err
		}

		y = given
	}

	return object.NewFloat(noise2D(x, y))
}

// mathRandom reads its range from how many arguments it is given: none gives a
// float in [0, 1), one gives a whole number in [1, n], and two give a whole
// number in [low, high].
func mathRandom(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.random", tok, args, 0, 2); err != nil {
		return err
	}

	if len(args) == 0 {
		return object.NewFloat(randomizer.Float64())
	}

	high, err := integerAt("math.random", tok, args, 0)

	if err != nil {
		return err
	}

	low := int64(1)

	if len(args) == 2 {
		bound, err := integerAt("math.random", tok, args, 1)

		if err != nil {
			return err
		}

		low = high
		high = bound
	}

	if high < low {
		return object.NewError(fault.Value, tok, "`math.random()` expects the upper bound to be at least the lower bound")
	}

	return object.NewInt(low + randomizer.Int63n(high-low+1))
}

// mathRandomSeed seeds the generator so that a run can be reproduced, which is
// what any procedurally generated result needs to be shareable. It seeds the
// same generator as random.seed().
func mathRandomSeed(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.randomSeed", tok, args, 1); err != nil {
		return err
	}

	given, err := integerAt("math.randomSeed", tok, args, 0)

	if err != nil {
		return err
	}

	SeedRandom(given)

	return value.NULL
}

// =============================================================================
// The operations themselves.

func numberSign(number *object.Number) *object.Number {
	if number.IsPos() {
		return object.NewInt(1)
	}

	if number.IsNeg() {
		return object.NewInt(-1)
	}

	return object.NewInt(0)
}

func numberSquare(number *object.Number) *object.Number {
	return number.Mul(number)
}

func numberReciprocal(tok token.Token, values []*object.Number) object.Object {
	if values[0].IsZero() {
		return divisionByZero("math.reciprocal", tok)
	}

	return object.NewFloat(1 / values[0].Float64())
}

func numberDivide(tok token.Token, values []*object.Number) object.Object {
	if values[1].IsZero() {
		return divisionByZero("math.divide", tok)
	}

	return values[0].Div(values[1])
}

func numberModulo(tok token.Token, values []*object.Number) object.Object {
	if values[1].IsZero() {
		return divisionByZero("math.mod", tok)
	}

	return values[0].Mod(values[1])
}

// numberPower keeps whole numbers whole where it can. An integer base raised to
// a non-negative integer power has an exact integer answer, and returning it as
// one means `math.pow(2, 10)` can index a list; when the answer would not fit,
// it falls back to floating point rather than wrapping.
func numberPower(tok token.Token, values []*object.Number) object.Object {
	base := values[0]
	exponent := values[1]

	if !base.IsFloat() && !exponent.IsFloat() && !exponent.IsNeg() {
		if result, ok := integerPower(base.Int64(), exponent.Int64()); ok {
			return object.NewInt(result)
		}
	}

	return object.NewFloat(math.Pow(base.Float64(), exponent.Float64()))
}

func numberMaximum(left *object.Number, right *object.Number) *object.Number {
	if right.GreaterThan(left) {
		return right
	}

	return left
}

func numberMinimum(left *object.Number, right *object.Number) *object.Number {
	if right.LessThan(left) {
		return right
	}

	return left
}

// numberClamp keeps a value inside a range. It answers with one of the three
// values it was given rather than a computed one, so clamping whole numbers
// leaves them whole.
func numberClamp(tok token.Token, values []*object.Number) object.Object {
	given := values[0]
	low := values[1]
	high := values[2]

	if low.GreaterThan(high) {
		return object.NewError(fault.Value, tok, "`math.clamp()` expects the lower bound to be no greater than the upper bound")
	}

	if given.LessThan(low) {
		return low
	}

	if given.GreaterThan(high) {
		return high
	}

	return given
}

// interpolate blends between two values. Smoothing, fades, and anything that
// should slide rather than jump is one of these per step.
func interpolate(from float64, to float64, amount float64) float64 {
	return from + (to-from)*amount
}

// smoothstep is interpolation that eases in and out, clamped to the edges.
func smoothstep(low float64, high float64, given float64) float64 {
	if low == high {
		if given < low {
			return 0
		}

		return 1
	}

	return smoothCurve(math.Max(0, math.Min(1, (given-low)/(high-low))))
}

func toDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func logGamma(given float64) float64 {
	result, _ := math.Lgamma(given)

	return result
}

func isNotANumber(number *object.Number) bool {
	return number.IsFloat() && math.IsNaN(number.Float64())
}

func isInfinite(number *object.Number) bool {
	return number.IsFloat() && math.IsInf(number.Float64(), 0)
}

func isFinite(number *object.Number) bool {
	return !isNotANumber(number) && !isInfinite(number)
}

func isInteger(number *object.Number) bool {
	if !number.IsFloat() {
		return true
	}

	given := number.Float64()

	return !math.IsNaN(given) && !math.IsInf(given, 0) && given == math.Trunc(given)
}

func isEven(number *object.Number) bool {
	return isInteger(number) && math.Mod(number.Float64(), 2) == 0
}

func isOdd(number *object.Number) bool {
	return isInteger(number) && math.Abs(math.Mod(number.Float64(), 2)) == 1
}

// closeEnough compares against both an absolute and a relative tolerance, so
// that it holds up for values near zero and for very large ones alike.
func closeEnough(left float64, right float64, tolerance float64) bool {
	if left == right {
		return true
	}

	difference := math.Abs(left - right)

	if difference <= tolerance {
		return true
	}

	return difference <= tolerance*math.Max(math.Abs(left), math.Abs(right))
}

// integerPower raises a whole number to a whole power, reporting false when the
// result would not fit so that the caller can fall back to floating point
// rather than silently wrapping.
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

// noise2D samples smoothed value noise at a point.
func noise2D(x float64, y float64) float64 {
	cellX := math.Floor(x)
	cellY := math.Floor(y)

	fractionX := smoothCurve(x - cellX)
	fractionY := smoothCurve(y - cellY)

	topLeft := latticeValue(int64(cellX), int64(cellY))
	topRight := latticeValue(int64(cellX)+1, int64(cellY))
	bottomLeft := latticeValue(int64(cellX), int64(cellY)+1)
	bottomRight := latticeValue(int64(cellX)+1, int64(cellY)+1)

	top := topLeft + (topRight-topLeft)*fractionX
	bottom := bottomLeft + (bottomRight-bottomLeft)*fractionX

	return top + (bottom-top)*fractionY
}

// smoothCurve eases an interpolation between 0 and 1 so that noise has no
// visible creases along cell boundaries.
func smoothCurve(given float64) float64 {
	return given * given * (3 - 2*given)
}

// latticeValue hashes integer coordinates into a repeatable value in [0, 1).
func latticeValue(x int64, y int64) float64 {
	hash := uint64(x)*0x9E3779B97F4A7C15 ^ uint64(y)*0xC2B2AE3D27D4EB4F

	hash ^= hash >> 33
	hash *= 0xFF51AFD7ED558CCD
	hash ^= hash >> 33

	return float64(hash>>11) / float64(uint64(1)<<53)
}

// =============================================================================
// Shared helpers

func toBoolean(given bool) object.Object {
	if given {
		return value.TRUE
	}

	return value.FALSE
}

func divisionByZero(name string, tok token.Token) *object.Error {
	return object.NewError(fault.Value, tok, "`%s()` cannot divide by zero", name)
}
