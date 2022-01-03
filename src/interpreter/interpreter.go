package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

type Evaluator func(node ast.Node, env *object.Environment) object.Object

// Evaluate is the heart of our interpreter. It switches off between the various
// AST node types and evaluates each accordingly and returns its value.
func Evaluate(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node, env)
	case *ast.Assign:
		return evaluateAssign(node, env)
	case *ast.Block:
		return evaluateBlock(node, env)
	case *ast.Boolean:
		return evaluateBoolean(node, env)
	case *ast.Call:
		return evaluateCall(node, env)
	case *ast.Expression:
		return Evaluate(node.Expression, env)
	case *ast.For:
		return evaluateFor(node, env)
	case *ast.Function:
		return evaluateFunction(node, env)
	case *ast.Identifier:
		return evaluateIdentifier(node, env)
	case *ast.If:
		return evaluateIf(node, env)
	case *ast.Index:
		return evaluateIndex(node, env)
	case *ast.Infix:
		return evaluateInfix(node, env)
	case *ast.List:
		return evaluateList(node, env)
	case *ast.Map:
		return evaluateMap(node, env)
	case *ast.Method:
		return evaluateMethod(node, env)
	case *ast.Null:
		return evaluateNull(node, env)
	case *ast.Number:
		return evaluateNumber(node, env)
	case *ast.Prefix:
		return evaluatePrefix(node, env)
	case *ast.String:
		return evaluateString(node, env)
	case *ast.While:
		return evaluateWhile(node, env)
	}

	return nil
}

// =============================================================================
// Helper functions

// toBooleanValue converts the passed native boolean value into a boolean
// object.
func toBooleanValue(input bool) *object.Boolean {
	if input {
		return value.TRUE
	}

	return value.FALSE
}

// isError determines if the referenced object is an error.
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}

	return false
}

// isTruthy determines if the referenced value is of a truthy value.
func isTruthy(value object.Object) bool {
	switch value := value.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return value.Value
	case *object.String:
		return len(value.Value) > 0
	default:
		return true
	}
}

// newError returns a new error object.
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
