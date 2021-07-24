package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
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
	default:
		return &object.Error{Message: "can only call functions and classes."}, false
	}
}
