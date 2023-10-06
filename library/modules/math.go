package modules

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
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

	return &object.Number{Value: number.Value.Abs()}
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

	return &object.Number{Value: number.Value.Cos()}
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

	return &object.Boolean{Value: number.Value.IsNegative()}
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

	return &object.Boolean{Value: number.Value.IsPositive()}
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

	return &object.Boolean{Value: number.Value.IsZero()}
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

	return &object.Number{Value: number.Value.Sin()}
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

	return &object.Number{Value: number.Value.Tan()}
}

// Properties

// mathPi returns the value of π, othewise known as Pi.
func mathPi(scope *object.Scope, tok token.Token) object.Object {
	pi, _ := decimal.NewFromString("3.141592653589793")

	return &object.Number{Value: pi}
}

// mathE returns the value of e, otherwise known as Euler's number.
func mathE(scope *object.Scope, tok token.Token) object.Object {
	e, _ := decimal.NewFromString("2.718281828459045")

	return &object.Number{Value: e}
}

// mathTau returns the value of τ, otherwise known as Tau. Tau is a circle
// constant equal to 2π, the ratio of a circle’s circumference to its radius.
func mathTau(scope *object.Scope, tok token.Token) object.Object {
	tau, _ := decimal.NewFromString("6.283185307179586")

	return &object.Number{Value: tau}
}

// mathEpsilon returns the value of ϵ, otherwise known as Epsilon. Epsilon
// represents the difference between 1 and the smallest floating point number
// greater than 1.
func mathEpsilon(scope *object.Scope, tok token.Token) object.Object {
	epsilon, _ := decimal.NewFromString("2.2204460492503130808472633361816E-16")

	return &object.Number{Value: epsilon}
}
