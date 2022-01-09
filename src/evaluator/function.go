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
		scope.Environment.Set(node.Name.Value, function)
	}

	return function
}

func createFunctionScope(function *object.Function, arguments []object.Object) *object.Scope {
	scope := &object.Scope{
		Self:        function,
		Environment: object.NewEnclosedEnvironment(function.Scope.Environment),
	}

	for key, val := range function.Defaults {
		scope.Environment.Set(key, Evaluate(val, scope))
	}

	for index, parameter := range function.Parameters {
		if index < len(arguments) {
			scope.Environment.Set(parameter.Value, arguments[index])
		}
	}

	return scope
}
