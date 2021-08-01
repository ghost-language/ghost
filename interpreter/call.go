package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateCall(node *ast.Call, env *object.Environment) (object.Object, bool) {
	callee, success := Evaluate(node.Callee, env)

	if !success {
		return &object.Error{Message: "unknown identifier."}, false
	}

	arguments := make([]object.Object, 0)

	for _, argument := range node.Arguments {
		nodeArgument, success := Evaluate(argument, env)

		if !success {
			return &object.Error{Message: "could not properly evaluate argument expressions."}, false
		}

		arguments = append(arguments, nodeArgument)
	}

	switch callable := callee.(type) {
	case *object.Class:
		instance := &object.ClassInstance{
			Class: callable,
		}

		return instance, true
	case *object.Standard:
		return callable.Function(arguments), true
	case *object.UserFunction:
		functionEnvironment := object.ExtendEnvironment(callable.Env)

		for i, parameter := range callable.Parameters {
			var val object.Object

			if i < len(arguments) {
				val = arguments[i]
			} else {
				val = value.NULL
			}

			functionEnvironment.Set(parameter.Name.Lexeme, val)
		}

		var evaluation object.Object
		var ok bool

		for _, statement := range callable.Body {
			evaluation, ok = Evaluate(statement, functionEnvironment)

			if !ok {
				return nil, ok
			}

			if returnValue, ok := evaluation.(*object.Return); ok {
				return returnValue.Value, ok
			}
		}

		return value.NULL, true
	default:
		return &object.Error{Message: fmt.Sprintf("can only call functions and classes, got=%s.", callable.String())}, false
	}
}