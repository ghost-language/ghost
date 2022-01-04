package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateFunction(node *ast.Function, env *object.Environment) object.Object {
	function := &object.Function{
		Parameters:  node.Parameters,
		Defaults:    node.Defaults,
		Body:        node.Body,
		Environment: env,
	}

	if node.Name != nil {
		env.Set(node.Name.Value, function)
	}

	return function
}

func createFunctionEnvironment(function *object.Function, arguments []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Environment)

	for key, val := range function.Defaults {
		env.Set(key, Evaluate(val, env))
	}

	for index, parameter := range function.Parameters {
		if index < len(arguments) {
			env.Set(parameter.Value, arguments[index])
		}
	}

	return env
}
