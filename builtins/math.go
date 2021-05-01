package builtins

import (
	"math/rand"
	"time"

	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	RegisterFunction("Math.abs", mathAbsFunction)
	RegisterFunction("Math.cos", mathCosFunction)
	RegisterFunction("Math.pi", mathPiFunction)
	RegisterFunction("Math.random", mathRandomFunction)
	RegisterFunction("Math.round", mathRoundFunction)
}

// mathAbsFunction returns the absolute value of the decimal.
func mathAbsFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() != object.NUMBER_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "Math.abs", "NUMBER", args[0].Type())
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Abs()}
}

// mathCosFunction returns the cosine of the radian decimal.
func mathCosFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if args[0].Type() != object.NUMBER_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "Math.cos", "NUMBER", args[0].Type())
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Cos()}
}

// mathPiFunction returns the value of pi. Will look into moving this to a property.
func mathPiFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	pi, _ := decimal.NewFromString("3.14159265358979323846264338327950288419716939937510582097494459")

	return &object.Number{Value: pi}
}

// mathRandomFunction returns a random decimal value with optional min max ranges.
func mathRandomFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	min := float64(0)
	max := float64(0)

	if len(args) > 0 {
		if args[0].Type() != object.NUMBER_OBJ {
			return error.NewError(line, error.ArgumentMustBe, "first", "Math.random", "NUMBER", args[0].Type())
		}

		max, _ = args[0].(*object.Number).Value.Float64()

		if len(args) > 1 {
			if args[1].Type() != object.NUMBER_OBJ {
				return error.NewError(line, error.ArgumentMustBe, "second", "Math.random", "NUMBER", args[1].Type())
			}

			min = max
			max, _ = args[1].(*object.Number).Value.Float64()

			if max < min {
				return error.NewError(line, error.ArgumentMustBe, "second", "Math.random", "larger than first argument", args[1].Type())
			}
		}
	}

	number := float64(0)

	if max > 0 {
		number = float64(min + rand.Float64() * (max - min))
	} else {
		number = rand.Float64()
	}

	return &object.Number{Value: decimal.NewFromFloat(number)}
}

func mathRoundFunction(env *object.Environment, line int, args ...object.Object) object.Object {
	if len(args) < 1 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 1)
	}

	if len(args) > 2 {
		return error.NewError(line, error.WrongNumberArguments, len(args), 2)
	}

	if args[0].Type() != object.NUMBER_OBJ {
		return error.NewError(line, error.ArgumentMustBe, "first", "Math.round", "NUMBER", args[0].Type())
	}

	places := &object.Number{Value: decimal.NewFromInt(0)}

	if len(args) == 2 {
		if args[1].Type() != object.NUMBER_OBJ {
			return error.NewError(line, error.ArgumentMustBe, "second", "Math.round", "NUMBER", args[1].Type())
		}

		places = args[1].(*object.Number)
	}

	number := args[0].(*object.Number)

	return &object.Number{Value: number.Value.Round(int32(places.Value.IntPart()))}
}