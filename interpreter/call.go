package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
)

func evaluateCall(node *ast.Call, env *environment.Environment) (object.Object, bool) {
	callee, success := Evaluate(node.Callee, env)

	if success != false {
		return nil, false
	}

	args := make([]object.Object, 0)

	for _, arg := range node.Arguments {
		nodeArg, success := Evaluate(arg, env)

		if success != true {
			return nil, false
		}

		args = append(args, nodeArg)
	}

	function, ok := callee.(*object.Standard)

	if !ok {
		return &object.Error{Message: fmt.Sprintf("can only call functions.")}, false
	}

	return function.Function(args), true
}
