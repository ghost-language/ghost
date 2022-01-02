package modules

import (
	"math/rand"
	"time"

	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

var Math = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Math, "abs", mathAbs)
	RegisterMethod(Math, "cos", mathCos)
	RegisterMethod(Math, "isNegative", mathIsNegative)
	RegisterMethod(Math, "isPositive", mathIsPositive)
	RegisterMethod(Math, "isZero", mathIsZero)
	RegisterMethod(Math, "pi", mathPi)
	RegisterMethod(Math, "random", mathRandom)
	RegisterMethod(Math, "seed", mathSeed)
	RegisterMethod(Math, "sin", mathSin)
	RegisterMethod(Math, "tan", mathTan)
}

// mathAbs returns the absolute value of the referenced number.
func mathAbs(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Abs()}
}

// mathCos returns the cosine value of the referenced number.
func mathCos(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Cos()}
}

// mathisNegative returns true if the referenced number is negative.
func mathIsNegative(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Boolean{Value: number.Value.IsNegative()}
}

// mathisPositive returns true if the referenced number is positive.
func mathIsPositive(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Boolean{Value: number.Value.IsPositive()}
}

// mathisZero returns true if the referenced number is zero.
func mathIsZero(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Boolean{Value: number.Value.IsZero()}
}

// mathPi returns the value of pi.
func mathPi(env *object.Environment, args ...object.Object) object.Object {
	pi, _ := decimal.NewFromString("3.14159265358979323846264338327950288419716939937510582097494459")

	return &object.Number{Value: pi}
}

// mathRandom when called without arguments returns a uniform pseudo-random real
// number in the range (0, 1). When called with a single number value (a), a
// pseudo-random number will be returned in the range (0, a). When called with
// two numbers (a, b), a pseudo-random number will be returned in the
// range (a, b).
func mathRandom(env *object.Environment, args ...object.Object) object.Object {
	min := float64(0)
	max := float64(0)

	if len(args) > 0 {
		max, _ = args[0].(*object.Number).Value.Float64()

		if len(args) > 1 {
			min = max
			max, _ = args[1].(*object.Number).Value.Float64()
		}
	}

	number := float64(0)

	if max > 0 {
		number = float64(min + rand.Float64()*(max-min))
	} else {
		number = rand.Float64()
	}

	return &object.Number{Value: decimal.NewFromFloat(number)}
}

// mathSeed sets the referenced number as the seed for the pseudo-random
// generator used by math.random(). If no value is passed, the current unix
// nano timestamp will be used.
func mathSeed(env *object.Environment, args ...object.Object) object.Object {
	var seed int64

	if len(args) == 1 && args[0].Type() == object.NUMBER {
		seed = args[0].(*object.Number).Value.IntPart()
	} else {
		seed = time.Now().UnixNano()
	}

	rand.Seed(seed)

	return nil
}

// mathSin returns the sine value of the referenced number.
func mathSin(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Sin()}
}

// mathTan returns the tangent value of the referenced number.
func mathTan(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Tan()}
}
