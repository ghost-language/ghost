package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateCall(node *ast.Call, env *object.Environment) (object.Object, bool) {
	callee, success := Evaluate(node.Callee, env)

	if !success {
		return &object.Error{Message: "unknown identifier."}, false
	}

	args := make([]object.Object, 0)

	for _, arg := range node.Arguments {
		nodeArg, success := Evaluate(arg, env)

		if !success {
			return &object.Error{Message: "could not properly evaluate argument expressions."}, false
		}

		args = append(args, nodeArg)
	}

	switch callable := callee.(type) {
	case *object.Class:
		instance := &object.ClassInstance{
			Class: callable,
		}

		return instance, true
	case *object.Standard:
		return callable.Function(args), true
	case *object.UserFunction:
		for _, parameter := range callable.Parameters {
			callable.Env.Set(parameter.Name.Lexeme, value.NULL)
		}

		for _, statement := range callable.Body {
			_, err := Evaluate(statement, callable.Env)

			if !err {
				return nil, err
			}
		}

		return nil, true
	default:
		return &object.Error{Message: "can only call functions and classes."}, false
	}
}
