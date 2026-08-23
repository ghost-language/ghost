package object

import (
	"io"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// Function objects consist of a user-generated function.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.Block
	Defaults   map[string]ast.ExpressionNode
	Scope      *Scope
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
func (function *Function) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}

// Evaluate evaluates the function's body ast.Block and returns the result.
func (function *Function) Evaluate(args []Object, writer io.Writer) Object {
	scope := function.scope(args)

	if writer != nil {
		scope.Environment.SetWriter(writer)
	}

	result := evaluator(function.Body, scope)

	return result
}

// Call invokes the function as a callback - the shape a library method wants
// when it takes a function argument, as List's map, filter, reduce, and each
// do. It unwraps the `return` marker Evaluate leaves in place, so a caller
// gets the value the Ghost function handed back rather than the object that
// carries it through evaluation.
func (function *Function) Call(args []Object) Object {
	result := function.Evaluate(args, nil)

	if wrapped, ok := result.(*Return); ok {
		return wrapped.Value
	}

	if result == nil {
		return &Null{}
	}

	return result
}

// =============================================================================
// Helper methods

func (function *Function) scope(arguments []Object) *Scope {
	scope := &Scope{
		Self:        function,
		Environment: NewEnclosedEnvironment(function.Scope.Environment),
	}

	for key, val := range function.Defaults {
		scope.Environment.Set(key, evaluator(val, scope))
	}

	for index, parameter := range function.Parameters {
		if index < len(arguments) {
			scope.Environment.Set(parameter.Value, arguments[index])
		}
	}

	return scope
}
