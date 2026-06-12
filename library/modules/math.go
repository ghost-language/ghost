package modules

import (
	"math"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var MathMethods = map[string]*object.LibraryFunction{}
var MathProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(MathMethods, "abs", mathAbs)
	RegisterMethod(MathMethods, "cos", mathCos)
	RegisterMethod(MathMethods, "isNegative", mathIsNegative)
	RegisterMethod(MathMethods, "isPositive", mathIsPositive)
	RegisterMethod(MathMethods, "isZero", mathIsZero)
	RegisterMethod(MathMethods, "sin", mathSin)
	RegisterMethod(MathMethods, "tan", mathTan)
	RegisterMethod(MathMethods, "max", mathMax)
	RegisterMethod(MathMethods, "min", mathMin)

	RegisterProperty(MathProperties, "pi", mathPi)
	RegisterProperty(MathProperties, "e", mathE)
	RegisterProperty(MathProperties, "epsilon", mathEpsilon)
	RegisterProperty(MathProperties, "tau", mathTau)
}

// mathAbs returns the absolute value of the referenced number.
func mathAbs(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return number.Abs()
}

// mathCos returns the cosine value of the referenced number.
func mathCos(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return number.Cos()
}

// mathisNegative returns true if the referenced number is negative.
func mathIsNegative(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Boolean{Value: number.IsNeg()}
}

// mathisPositive returns true if the referenced number is positive.
func mathIsPositive(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Boolean{Value: number.IsPos()}
}

// mathisZero returns true if the referenced number is zero.
func mathIsZero(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return &object.Boolean{Value: number.IsZero()}
}

// mathSin returns the sine value of the referenced number.
func mathSin(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return number.Sin()
}

// mathTan returns the tangent value of the referenced number.
func mathTan(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	number := args[0].(*object.Number)

	return number.Tan()
}

// mathMax returns the largest number of the referenced numbers.
func mathMax(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError("%d:%d:%s: runtime error: math.max requires at least two arguments", tok.Line, tok.Column, tok.File)
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	if args[1].Type() != object.NUMBER {
		return nil
	}

	number1 := args[0].(*object.Number)
	number2 := args[1].(*object.Number)

	if number1.GreaterThan(number2) {
		return number1
	}

	return number2
}

// mathMin returns the smallest number of the referenced numbers.
func mathMin(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError("%d:%d:%s: runtime error: math.min requires at least two arguments", tok.Line, tok.Column, tok.File)
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	if args[1].Type() != object.NUMBER {
		return nil
	}

	number1 := args[0].(*object.Number)
	number2 := args[1].(*object.Number)

	if number1.LessThan(number2) {
		return number1
	}

	return number2
}

// Properties

// mathPi returns the value of π, othewise known as Pi.
func mathPi(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(math.Pi)
}

// mathE returns the value of e, otherwise known as Euler's number.
func mathE(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(math.E)
}

// mathTau returns the value of τ, otherwise known as Tau. Tau is a circle
// constant equal to 2π, the ratio of a circle's circumference to its radius.
func mathTau(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(2 * math.Pi)
}

// mathEpsilon returns the value of ϵ, otherwise known as Epsilon. Epsilon
// represents the difference between 1 and the smallest floating point number
// greater than 1.
func mathEpsilon(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(math.SmallestNonzeroFloat64)
}
