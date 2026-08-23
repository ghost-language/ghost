package modules

import (
	"math"
	"testing"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// call invokes a registered math method the way the evaluator would.
func call(t *testing.T, name string, args ...object.Object) object.Object {
	t.Helper()

	method, ok := MathMethods[name]

	if !ok {
		t.Fatalf("math.%s is not registered", name)
	}

	return method.Function(nil, token.Token{}, args...)
}

func property(t *testing.T, name string) object.Object {
	t.Helper()

	constant, ok := MathProperties[name]

	if !ok {
		t.Fatalf("math.%s is not registered", name)
	}

	return constant.Property(nil, token.Token{})
}

func list(values ...object.Object) *object.List {
	return &object.List{Elements: values}
}

func integers(values ...int64) *object.List {
	return integerList(values)
}

// stringOf renders a result the way a Ghost program would see it, which keeps
// the expectations below readable for lists and matrices alike.
func stringOf(t *testing.T, result object.Object) string {
	t.Helper()

	if result == nil {
		t.Fatalf("method returned nothing")
	}

	if object.IsError(result) {
		t.Fatalf("unexpected error: %s", result.String())
	}

	return result.String()
}

func TestMathScalars(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Object
		expected string
	}{
		{"abs", []object.Object{object.NewInt(-5)}, "5"},
		{"sign", []object.Object{object.NewFloat(-0.5)}, "-1"},
		{"sqrt", []object.Object{object.NewInt(16)}, "4"},
		{"cbrt", []object.Object{object.NewInt(27)}, "3"},
		{"square", []object.Object{object.NewInt(7)}, "49"},
		{"reciprocal", []object.Object{object.NewInt(4)}, "0.25"},
		{"floor", []object.Object{object.NewFloat(2.7)}, "2"},
		{"ceil", []object.Object{object.NewFloat(2.1)}, "3"},
		{"truncate", []object.Object{object.NewFloat(-2.7)}, "-2"},
		{"round", []object.Object{object.NewFloat(2.5)}, "3"},
		{"round", []object.Object{object.NewFloat(2.567), object.NewInt(2)}, "2.57"},
		{"pow", []object.Object{object.NewInt(2), object.NewInt(10)}, "1024"},
		{"hypot", []object.Object{object.NewInt(3), object.NewInt(4)}, "5"},
		{"log", []object.Object{object.NewInt(8), object.NewInt(2)}, "3"},
		{"log2", []object.Object{object.NewInt(8)}, "3"},
		{"log10", []object.Object{object.NewInt(1000)}, "3"},
		{"degrees", []object.Object{object.NewFloat(math.Pi)}, "180"},
		{"radians", []object.Object{object.NewInt(180)}, "3.141592653589793"},
		{"atan2", []object.Object{object.NewInt(1), object.NewInt(1)}, "0.7853981633974483"},
		{"clamp", []object.Object{object.NewInt(15), object.NewInt(0), object.NewInt(10)}, "10"},
		{"lerp", []object.Object{object.NewInt(0), object.NewInt(100), object.NewFloat(0.25)}, "25"},
		{"smoothstep", []object.Object{object.NewInt(0), object.NewInt(1), object.NewFloat(0.5)}, "0.5"},
		{"gamma", []object.Object{object.NewInt(5)}, "24"},
		{"mod", []object.Object{object.NewInt(10), object.NewInt(3)}, "1"},
		{"divide", []object.Object{object.NewInt(10), object.NewInt(4)}, "2.5"},
		{"maximum", []object.Object{object.NewInt(3), object.NewInt(7)}, "7"},
		{"minimum", []object.Object{object.NewInt(3), object.NewInt(7)}, "3"},
	}

	for _, test := range tests {
		if result := stringOf(t, call(t, test.name, test.args...)); result != test.expected {
			t.Errorf("math.%s: got=%s, expected=%s", test.name, result, test.expected)
		}
	}
}

// TestMathPreservesWholeNumbers pins the rule that an exact answer to a whole
// number question comes back as a whole number, so that results can index a
// list without a conversion.
func TestMathPreservesWholeNumbers(t *testing.T) {
	tests := []struct {
		name string
		args []object.Object
	}{
		{"abs", []object.Object{object.NewInt(-5)}},
		{"square", []object.Object{object.NewInt(7)}},
		{"pow", []object.Object{object.NewInt(2), object.NewInt(10)}},
		{"floor", []object.Object{object.NewFloat(2.7)}},
		{"add", []object.Object{object.NewInt(2), object.NewInt(3)}},
		{"clamp", []object.Object{object.NewInt(15), object.NewInt(0), object.NewInt(10)}},
		{"sum", []object.Object{integers(1, 2, 3)}},
		{"max", []object.Object{integers(1, 2, 3)}},
		{"factorial", []object.Object{object.NewInt(10)}},
	}

	for _, test := range tests {
		result := call(t, test.name, test.args...)
		number, ok := result.(*object.Number)

		if !ok {
			t.Fatalf("math.%s: got=%s, expected a number", test.name, stringOf(t, result))
		}

		if number.IsFloat() {
			t.Errorf("math.%s: got a float, expected a whole number", test.name)
		}
	}
}

// TestMathPowerOverflowsToFloat covers the fallback that keeps an answer too
// large for a whole number from wrapping around.
func TestMathPowerOverflowsToFloat(t *testing.T) {
	result := call(t, "pow", object.NewInt(2), object.NewInt(200))
	number, ok := result.(*object.Number)

	if !ok || !number.IsFloat() {
		t.Fatalf("expected a float, got=%s", stringOf(t, result))
	}

	if !closeEnough(number.Float64(), math.Pow(2, 200), 1e-9) {
		t.Errorf("got=%s, expected=%v", result.String(), math.Pow(2, 200))
	}
}

// TestMathBroadcasts is the heart of the module: an operation written for one
// number works on a list, and on a list of lists, without knowing it.
func TestMathBroadcasts(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Object
		expected string
	}{
		{"sqrt", []object.Object{integers(1, 4, 9)}, "[1, 2, 3]"},
		{"sqrt", []object.Object{list(integers(1, 4), integers(9, 16))}, "[[1, 2], [3, 4]]"},
		{"add", []object.Object{integers(1, 2, 3), object.NewInt(10)}, "[11, 12, 13]"},
		{"add", []object.Object{object.NewInt(10), integers(1, 2, 3)}, "[11, 12, 13]"},
		{"multiply", []object.Object{integers(1, 2, 3), integers(4, 5, 6)}, "[4, 10, 18]"},
		{"multiply", []object.Object{list(integers(1, 2)), object.NewInt(3)}, "[[3, 6]]"},
		{"round", []object.Object{list(object.NewFloat(1.24), object.NewFloat(1.26)), object.NewInt(1)}, "[1.2, 1.3]"},
		{"clamp", []object.Object{integers(-5, 5, 15), object.NewInt(0), object.NewInt(10)}, "[0, 5, 10]"},
		{"isNegative", []object.Object{integers(-1, 0, 1)}, "[true, false, false]"},
		{"maximum", []object.Object{integers(1, 5), integers(4, 2)}, "[4, 5]"},

		// Shapes are lined up from the right and the shorter side stretches,
		// which is numpy's rule. A row spreads down a matrix rather than being
		// paired against its rows.
		{"add", []object.Object{list(integers(1, 2), integers(3, 4)), integers(10, 20)}, "[[11, 22], [13, 24]]"},
		{"add", []object.Object{list(integers(1, 2), integers(3, 4)), list(integers(10), integers(20))}, "[[11, 12], [23, 24]]"},
		{"add", []object.Object{integers(1, 2, 3), integers(5)}, "[6, 7, 8]"},
		{"clamp", []object.Object{list(integers(-5, 5), integers(15, 3)), object.NewInt(0), object.NewInt(10)}, "[[0, 5], [10, 3]]"},
	}

	for _, test := range tests {
		if result := stringOf(t, call(t, test.name, test.args...)); result != test.expected {
			t.Errorf("math.%s: got=%s, expected=%s", test.name, result, test.expected)
		}
	}
}

// TestMathReductionsAcceptEveryShape pins that a reduction reads its values
// spread across the call, collected in a list, or arranged as a matrix.
func TestMathReductionsAcceptEveryShape(t *testing.T) {
	shapes := [][]object.Object{
		{object.NewInt(1), object.NewInt(2), object.NewInt(3), object.NewInt(4)},
		{integers(1, 2, 3, 4)},
		{list(integers(1, 2), integers(3, 4))},
	}

	for _, args := range shapes {
		if result := stringOf(t, call(t, "sum", args...)); result != "10" {
			t.Errorf("math.sum: got=%s, expected=10", result)
		}
	}
}

func TestMathStatistics(t *testing.T) {
	sample := integers(2, 4, 4, 4, 5, 5, 7, 9)

	tests := []struct {
		name     string
		args     []object.Object
		expected string
	}{
		{"sum", []object.Object{integers(1, 2, 3, 4)}, "10"},
		{"product", []object.Object{integers(1, 2, 3, 4)}, "24"},
		{"mean", []object.Object{sample}, "5"},
		{"median", []object.Object{integers(1, 3, 2)}, "2"},
		{"median", []object.Object{integers(1, 2, 3, 4)}, "2.5"},
		{"mode", []object.Object{integers(1, 2, 2, 3)}, "2"},
		{"variance", []object.Object{sample}, "4"},
		{"standardDeviation", []object.Object{sample}, "2"},
		{"sampleVariance", []object.Object{integers(2, 4, 4, 4, 5, 5, 7, 9)}, "4.571428571428571"},
		{"min", []object.Object{integers(3, 1, 2)}, "1"},
		{"max", []object.Object{integers(3, 1, 2)}, "3"},
		{"argmin", []object.Object{integers(3, 1, 2)}, "1"},
		{"argmax", []object.Object{integers(3, 1, 2)}, "0"},
		{"percentile", []object.Object{integers(1, 2, 3, 4), object.NewInt(50)}, "2.5"},
		{"quantile", []object.Object{integers(1, 2, 3, 4), object.NewFloat(0.25)}, "1.75"},
		{"cumulativeSum", []object.Object{integers(1, 2, 3)}, "[1, 3, 6]"},
		{"cumulativeProduct", []object.Object{integers(1, 2, 3, 4)}, "[1, 2, 6, 24]"},
		{"sort", []object.Object{integers(3, 1, 2)}, "[1, 2, 3]"},
		{"sort", []object.Object{integers(3, 1, 2), &object.Boolean{Value: true}}, "[3, 2, 1]"},
		{"unique", []object.Object{integers(1, 1, 2, 3, 3)}, "[1, 2, 3]"},
		{"gcd", []object.Object{integers(12, 18, 24)}, "6"},
		{"lcm", []object.Object{integers(4, 6)}, "12"},
		{"factorial", []object.Object{object.NewInt(10)}, "3628800"},
		{"combinations", []object.Object{object.NewInt(5), object.NewInt(2)}, "10"},
		{"permutations", []object.Object{object.NewInt(5), object.NewInt(2)}, "20"},
	}

	for _, test := range tests {
		if result := stringOf(t, call(t, test.name, test.args...)); result != test.expected {
			t.Errorf("math.%s: got=%s, expected=%s", test.name, result, test.expected)
		}
	}
}

func TestMathArrays(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Object
		expected string
	}{
		{"arange", []object.Object{object.NewInt(5)}, "[0, 1, 2, 3, 4]"},
		{"arange", []object.Object{object.NewInt(1), object.NewInt(4)}, "[1, 2, 3]"},
		{"arange", []object.Object{object.NewInt(0), object.NewInt(1), object.NewFloat(0.25)}, "[0, 0.25, 0.5, 0.75]"},
		{"linspace", []object.Object{object.NewInt(0), object.NewInt(1), object.NewInt(5)}, "[0, 0.25, 0.5, 0.75, 1]"},
		{"zeros", []object.Object{object.NewInt(3)}, "[0, 0, 0]"},
		{"ones", []object.Object{object.NewInt(2), object.NewInt(3)}, "[[1, 1, 1], [1, 1, 1]]"},
		{"full", []object.Object{object.NewInt(2), object.NewInt(7)}, "[7, 7]"},
		{"identity", []object.Object{object.NewInt(2)}, "[[1, 0], [0, 1]]"},
		{"reshape", []object.Object{integers(1, 2, 3, 4, 5, 6), object.NewInt(2), object.NewInt(3)}, "[[1, 2, 3], [4, 5, 6]]"},
		{"reshape", []object.Object{integers(1, 2, 3, 4, 5, 6), object.NewInt(-1), object.NewInt(2)}, "[[1, 2], [3, 4], [5, 6]]"},
		{"flatten", []object.Object{list(integers(1, 2), integers(3, 4))}, "[1, 2, 3, 4]"},
		{"shape", []object.Object{list(integers(1, 2, 3), integers(4, 5, 6))}, "[2, 3]"},
		{"shape", []object.Object{integers(1, 2, 3)}, "[3]"},
		{"transpose", []object.Object{list(integers(1, 2, 3), integers(4, 5, 6))}, "[[1, 4], [2, 5], [3, 6]]"},
	}

	for _, test := range tests {
		if result := stringOf(t, call(t, test.name, test.args...)); result != test.expected {
			t.Errorf("math.%s: got=%s, expected=%s", test.name, result, test.expected)
		}
	}
}

func TestMathLinearAlgebra(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Object
		expected string
	}{
		{"dot", []object.Object{integers(1, 2, 3), integers(4, 5, 6)}, "32"},
		{"dot", []object.Object{list(integers(1, 2), integers(3, 4)), list(integers(5, 6), integers(7, 8))}, "[[19, 22], [43, 50]]"},
		{"dot", []object.Object{list(integers(1, 2), integers(3, 4)), integers(5, 6)}, "[17, 39]"},
		{"matmul", []object.Object{list(integers(1, 2)), list(integers(3), integers(4))}, "[[11]]"},
		{"cross", []object.Object{integers(1, 0, 0), integers(0, 1, 0)}, "[0, 0, 1]"},
		{"cross", []object.Object{integers(1, 0), integers(0, 1)}, "1"},
		{"outer", []object.Object{integers(1, 2, 3), integers(10, 20)}, "[[10, 20], [20, 40], [30, 60]]"},
		{"norm", []object.Object{integers(3, 4)}, "5"},
		{"norm", []object.Object{integers(3, -4), object.NewInt(1)}, "7"},
		{"normalize", []object.Object{integers(3, 4)}, "[0.6, 0.8]"},
		{"determinant", []object.Object{list(integers(1, 2), integers(3, 4))}, "-2"},
		{"trace", []object.Object{list(integers(1, 2), integers(3, 4))}, "5"},
		{"solve", []object.Object{list(integers(2, 1), integers(1, 3)), integers(5, 10)}, "[1, 3]"},
		{"distance", []object.Object{object.NewInt(0), object.NewInt(0), object.NewInt(3), object.NewInt(4)}, "5"},
		{"distance", []object.Object{integers(0, 0), integers(3, 4)}, "5"},
		{"distance", []object.Object{integers(0, 0, 0), integers(1, 2, 2)}, "3"},
		{"angle", []object.Object{object.NewInt(0), object.NewInt(0), object.NewInt(1), object.NewInt(1)}, "0.7853981633974483"},
	}

	for _, test := range tests {
		if result := stringOf(t, call(t, test.name, test.args...)); result != test.expected {
			t.Errorf("math.%s: got=%s, expected=%s", test.name, result, test.expected)
		}
	}
}

// TestMathInverseUndoesMultiplication checks the inverse against the property
// that defines it, rather than against digits that depend on rounding.
func TestMathInverseUndoesMultiplication(t *testing.T) {
	matrix := list(integers(4, 7), integers(2, 6))
	inverse := call(t, "inverse", matrix)

	if object.IsError(inverse) {
		t.Fatalf("unexpected error: %s", inverse.String())
	}

	product := call(t, "dot", matrix, inverse)
	rows, _, err := toMatrix("test", token.Token{}, []object.Object{product}, 0)

	if err != nil {
		t.Fatalf("unexpected error: %s", err.String())
	}

	for row := range rows {
		for column := range rows[row] {
			expected := 0.0

			if row == column {
				expected = 1.0
			}

			if !closeEnough(rows[row][column], expected, 1e-9) {
				t.Errorf("inverse product at %d,%d: got=%v, expected=%v", row, column, rows[row][column], expected)
			}
		}
	}
}

func TestMathPredicates(t *testing.T) {
	tests := []struct {
		name     string
		args     []object.Object
		expected bool
	}{
		{"isNaN", []object.Object{object.NewFloat(math.NaN())}, true},
		{"isNaN", []object.Object{object.NewInt(1)}, false},
		{"isFinite", []object.Object{object.NewInt(1)}, true},
		{"isInfinite", []object.Object{object.NewFloat(math.Inf(1))}, true},
		{"isInteger", []object.Object{object.NewFloat(2.0)}, true},
		{"isInteger", []object.Object{object.NewFloat(2.5)}, false},
		{"isEven", []object.Object{object.NewInt(4)}, true},
		{"isOdd", []object.Object{object.NewInt(3)}, true},
		{"isOdd", []object.Object{object.NewInt(-3)}, true},
		{"isPrime", []object.Object{object.NewInt(97)}, true},
		{"isPrime", []object.Object{object.NewInt(1)}, false},
		{"isZero", []object.Object{object.NewInt(0)}, true},
		{"isClose", []object.Object{object.NewFloat(0.1 + 0.2), object.NewFloat(0.3)}, true},
		{"isClose", []object.Object{object.NewFloat(0.1), object.NewFloat(0.3)}, false},
		{"isClose", []object.Object{object.NewFloat(0.1), object.NewFloat(0.3), object.NewFloat(0.5)}, true},
	}

	for _, test := range tests {
		result := call(t, test.name, test.args...)
		boolean, ok := result.(*object.Boolean)

		if !ok {
			t.Fatalf("math.%s: got=%s, expected a boolean", test.name, stringOf(t, result))
		}

		if boolean.Value != test.expected {
			t.Errorf("math.%s: got=%t, expected=%t", test.name, boolean.Value, test.expected)
		}
	}
}

func TestMathConstants(t *testing.T) {
	tests := []struct {
		name     string
		expected float64
	}{
		{"pi", math.Pi},
		{"tau", 2 * math.Pi},
		{"e", math.E},
		{"phi", math.Phi},
		{"sqrt2", math.Sqrt2},
		{"ln2", math.Ln2},
		{"epsilon", math.Nextafter(1, 2) - 1},
		{"smallestNumber", math.SmallestNonzeroFloat64},
		{"largestNumber", math.MaxFloat64},
	}

	for _, test := range tests {
		result := property(t, test.name)
		number, ok := result.(*object.Number)

		if !ok {
			t.Fatalf("math.%s: expected a number", test.name)
		}

		if number.Float64() != test.expected {
			t.Errorf("math.%s: got=%v, expected=%v", test.name, number.Float64(), test.expected)
		}
	}
}

// TestMathRandomIsReproducible covers the seeding contract: the same seed
// replays the same numbers, and math.randomSeed and random.seed drive the same
// generator.
func TestMathRandomIsReproducible(t *testing.T) {
	call(t, "randomSeed", object.NewInt(42))
	first := stringOf(t, call(t, "randomInt", object.NewInt(1), object.NewInt(1000)))

	call(t, "randomSeed", object.NewInt(42))
	second := stringOf(t, call(t, "randomInt", object.NewInt(1), object.NewInt(1000)))

	if first != second {
		t.Errorf("seeded runs differ. got=%s and %s", first, second)
	}

	randomSeed(nil, token.Token{}, object.NewInt(42))
	third := stringOf(t, call(t, "randomInt", object.NewInt(1), object.NewInt(1000)))

	if third != first {
		t.Errorf("random.seed() and math.randomSeed() drive different generators. got=%s and %s", third, first)
	}
}

// TestMathNoiseIsSmoothAndRepeatable covers the two properties noise is chosen
// for: the same input always gives the same output, and nearby inputs give
// nearby outputs.
func TestMathNoiseIsSmoothAndRepeatable(t *testing.T) {
	first := call(t, "noise", object.NewFloat(3.7), object.NewFloat(1.2))
	second := call(t, "noise", object.NewFloat(3.7), object.NewFloat(1.2))

	if stringOf(t, first) != stringOf(t, second) {
		t.Errorf("noise is not repeatable. got=%s and %s", first.String(), second.String())
	}

	for step := 0; step < 100; step++ {
		x := float64(step) / 10

		value := noise2D(x, 0)
		next := noise2D(x+0.01, 0)

		if value < 0 || value > 1 {
			t.Fatalf("noise left the 0-1 range at %v. got=%v", x, value)
		}

		if math.Abs(next-value) > 0.1 {
			t.Errorf("noise jumped between %v and %v. got=%v and %v", x, x+0.01, value, next)
		}
	}
}

func TestMathErrors(t *testing.T) {
	tests := []struct {
		name string
		args []object.Object
	}{
		{"sqrt", []object.Object{&object.String{Value: "four"}}},
		{"sqrt", []object.Object{}},
		{"add", []object.Object{integers(1, 2), integers(1, 2, 3)}},
		{"divide", []object.Object{object.NewInt(1), object.NewInt(0)}},
		{"mod", []object.Object{object.NewInt(1), object.NewInt(0)}},
		{"reciprocal", []object.Object{object.NewInt(0)}},
		{"clamp", []object.Object{object.NewInt(1), object.NewInt(10), object.NewInt(0)}},
		{"factorial", []object.Object{object.NewInt(-1)}},
		{"combinations", []object.Object{object.NewInt(2), object.NewInt(5)}},
		{"percentile", []object.Object{integers(1, 2), object.NewInt(150)}},
		{"sampleVariance", []object.Object{integers(1)}},
		{"normalize", []object.Object{integers(0, 0)}},
		{"inverse", []object.Object{list(integers(1, 2), integers(2, 4))}},
		{"solve", []object.Object{list(integers(1, 2), integers(2, 4)), integers(1, 2)}},
		{"determinant", []object.Object{list(integers(1, 2, 3), integers(4, 5, 6))}},
		{"dot", []object.Object{integers(1, 2), integers(1, 2, 3)}},
		{"outer", []object.Object{integers(1, 2), object.NewInt(3)}},
		{"reshape", []object.Object{integers(1, 2, 3), object.NewInt(2), object.NewInt(2)}},
		{"arange", []object.Object{object.NewInt(1), object.NewInt(5), object.NewInt(0)}},
		{"linspace", []object.Object{object.NewInt(0), object.NewInt(1), object.NewInt(-1)}},
		{"identity", []object.Object{object.NewInt(-1)}},
		{"randomInt", []object.Object{object.NewInt(10), object.NewInt(1)}},
		{"max", []object.Object{}},
		{"transpose", []object.Object{object.NewInt(5)}},
	}

	for _, test := range tests {
		result := call(t, test.name, test.args...)

		if !object.IsError(result) {
			t.Errorf("math.%s: expected an error, got=%v", test.name, result)
		}
	}
}
