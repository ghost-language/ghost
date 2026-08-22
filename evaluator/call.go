package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

// maxCallDepth bounds how deeply calls may nest.
//
// Ghost runs on the Go stack, and a Go stack overflow is not something a
// program can recover from: it kills the process outright, with a Go traceback
// and no mention of the Ghost code that caused it. A recursion that never
// bottoms out is an ordinary mistake, so the depth is counted and reported as
// an ordinary error while there is still a stack to report it on.
const maxCallDepth = 4096

func evaluateCall(node *ast.Call, scope *object.Scope) object.Object {
	callee := Evaluate(node.Callee, scope)

	if isError(callee) {
		return callee
	}

	arguments := evaluateExpressions(node.Arguments, scope)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	// The call is reported at the thing being called rather than at the bracket
	// that follows it, which is where a reader looks when told a call went wrong.
	at := calleeToken(node.Callee, node.Token)

	result := unwrapCall(at, callee, arguments, scope)

	// The frame is recorded as the error passes back through the call, so a
	// failure deep inside a helper is reported with the path that reached it.
	// Only calls into Ghost code get one: a library function reports its own
	// failures at the call itself, and a frame there would just repeat the
	// position already printed above it.
	if failed, ok := result.(*object.Error); ok {
		if _, written := callee.(*object.Function); written {
			return failed.WithFrame(callName(node.Callee), at)
		}
	}

	return result
}

// calleeToken finds the token that names what is being called, falling back to
// the call itself for an expression with no name of its own.
func calleeToken(callee ast.ExpressionNode, fallback token.Token) token.Token {
	switch callee := callee.(type) {
	case *ast.Identifier:
		return callee.Token
	case *ast.Property:
		if property, ok := callee.Property.(*ast.Identifier); ok {
			return property.Token
		}
	}

	return fallback
}

// callName reads the name a call was written with, for the stack trace. An
// expression that is not simply a name — a function pulled out of a list, say —
// has no name to give, and the frame says so rather than guessing.
func callName(callee ast.ExpressionNode) string {
	switch callee := callee.(type) {
	case *ast.Identifier:
		return callee.Value + "()"
	case *ast.Property:
		if property, ok := callee.Property.(*ast.Identifier); ok {
			return property.Value + "()"
		}
	}

	return ""
}

func unwrapCall(tok token.Token, callee object.Object, arguments []object.Object, scope *object.Scope) object.Object {
	if callee == nil {
		return object.NewError(fault.Type, tok, "cannot call a null value")
	}

	switch callee := callee.(type) {
	case *object.LibraryFunction:
		if result := callee.Function(scope, tok, arguments...); result != nil {
			return result
		}

		return nil
	case *object.LibraryProperty:
		if result := callee.Property(scope, tok); result != nil {
			return result
		}

		return nil
	case *object.Function:
		if scope.Depth >= maxCallDepth {
			return tooDeep(tok)
		}

		functionEnvironment := createFunctionEnvironment(callee, arguments)
		functionScope := &object.Scope{Self: callee, Environment: functionEnvironment, Depth: scope.Depth + 1}

		// A function declared inside a method keeps that method's receiver, so
		// `this` and `super` still work from within a closure.
		if callee.Scope != nil {
			if instance, ok := callee.Scope.Self.(*object.Instance); ok {
				functionScope.Self = instance
				functionScope.Class = callee.Scope.Class
			}
		}

		evaluated := Evaluate(callee.Body, functionScope)

		return unwrapReturn(evaluated)
	default:
		return object.NewError(fault.Type, tok, "cannot call %s, which is not a function", object.TypeName(callee))
	}
}

func unwrapReturn(obj object.Object) object.Object {
	switch value := obj.(type) {
	case *object.Error:
		return obj
	case *object.Return:
		return value.Value
	}

	return value.NULL
}
