package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateFunction(node *ast.Function, scope *object.Scope) object.Object {
	function := &object.Function{
		Parameters: node.Parameters,
		Defaults:   node.Defaults,
		Body:       node.Body,
		Scope:      scope,
	}

	if node.Name != nil {
		switch this := scope.Self.(type) {
		case *object.Class:
			this.Environment.Set(node.Name.Value, function)
		default:
			scope.Environment.Set(node.Name.Value, function)
		}
	}

	return function
}

func createFunctionEnvironment(function *object.Function, arguments []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Scope.Environment)

	for key, val := range function.Defaults {
		env.Set(key, Evaluate(val, function.Scope))
	}

	for index, parameter := range function.Parameters {
		if index < len(arguments) {
			env.Set(parameter.Value, arguments[index])
		}
	}

	return env
}
