package object

import (
	"io"

	"ghostlang.org/x/ghost/ast"
)

const FUNCTION = "FUNCTION"

// Function objects consist of a user-generated function.
type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.Block
	Defaults    map[string]ast.ExpressionNode
	Environment *Environment
}

// String represents the function object's value as a string.
func (function *Function) String() string {
	return "function"
}

// Type returns the function object type.
func (function *Function) Type() Type {
	return FUNCTION
}

// Method defines the set of methods available on function objects.
func (function *Function) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

// Evaluate evaluates the function's body ast.Block and returns the result.
func (function *Function) Evaluate(args []Object, writer io.Writer) Object {
	env := function.environment(args)

	if writer != nil {
		env.SetWriter(writer)
	}

	result := evaluator(function.Body, env)

	return result
}

// =============================================================================
// Helper methods

func (function *Function) environment(arguments []Object) *Environment {
	env := NewEnclosedEnvironment(function.Environment)

	for key, val := range function.Defaults {
		result := evaluator(val, env)
		env.Set(key, result)
	}

	for index, parameter := range function.Parameters {
		if index < len(arguments) {
			env.Set(parameter.Value, arguments[index])
		}
	}

	return env
}
