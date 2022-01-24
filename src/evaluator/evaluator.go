package evaluator

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

type Evaluator func(node ast.Node, scope *object.Scope) object.Object

// Evaluate is the heart of our evaluator. It switches off between the various
// AST node types and evaluates each accordingly and returns its value.
func Evaluate(node ast.Node, scope *object.Scope) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node, scope)
	case *ast.Assign:
		return evaluateAssign(node, scope)
	case *ast.Block:
		return evaluateBlock(node, scope)
	case *ast.Boolean:
		return evaluateBoolean(node, scope)
	case *ast.Call:
		return evaluateCall(node, scope)
	case *ast.Class:
		return evaluateClass(node, scope)
	case *ast.Compound:
		return evaluateCompound(node, scope)
	case *ast.Expression:
		return Evaluate(node.Expression, scope)
	case *ast.For:
		return evaluateFor(node, scope)
	case *ast.ForIn:
		return evaluateForIn(node, scope)
	case *ast.Function:
		return evaluateFunction(node, scope)
	case *ast.Identifier:
		return evaluateIdentifier(node, scope)
	case *ast.If:
		return evaluateIf(node, scope)
	case *ast.Import:
		return evaluateImport(node, scope)
	case *ast.ImportFrom:
		return evaluateImportFrom(node, scope)
	case *ast.Index:
		return evaluateIndex(node, scope)
	case *ast.Infix:
		return evaluateInfix(node, scope)
	case *ast.List:
		return evaluateList(node, scope)
	case *ast.Map:
		return evaluateMap(node, scope)
	case *ast.Method:
		return evaluateMethod(node, scope)
	case *ast.Null:
		return evaluateNull(node, scope)
	case *ast.Number:
		return evaluateNumber(node, scope)
	case *ast.Prefix:
		return evaluatePrefix(node, scope)
	case *ast.Property:
		return evaluateProperty(node, scope)
	case *ast.Return:
		return evaluateReturn(node, scope)
	case *ast.String:
		return evaluateString(node, scope)
	case *ast.Switch:
		return evaluateSwitch(node, scope)
	case *ast.This:
		return evaluateThis(node, scope)
	case *ast.While:
		return evaluateWhile(node, scope)
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
