package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateFunction(node *ast.Function, scope *object.Scope) object.Object {
	function := &object.Function{
		Parameters: node.Parameters,
		Defaults:   node.Defaults,
		Body:       node.Body,
		Scope:      scope,
		Rest:       node.Rest,
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

// createFunctionEnvironment binds arguments to a user-defined function's
// parameters. name and tok name and locate the call for an arity error - see
// checkArity - the same way every library call already reports one.
func createFunctionEnvironment(function *object.Function, arguments []object.Object, name string, tok token.Token) (*object.Environment, *object.Error) {
	if err := checkArity(function, arguments, name, tok); err != nil {
		return nil, err
	}

	env := object.NewEnclosedEnvironment(function.Scope.Environment)

	for key, val := range function.Defaults {
		env.Set(key, Evaluate(val, function.Scope))
	}

	fixed := len(function.Parameters)

	if function.Rest {
		fixed--
	}

	for index := 0; index < fixed; index++ {
		if index < len(arguments) {
			env.Set(function.Parameters[index].Value, arguments[index])
		}
	}

	// A rest parameter collects everything from its position onward into a
	// list of its own - always a list, empty rather than absent when nothing
	// was left to collect, so the parameter never needs a nil check.
	if function.Rest {
		rest := []object.Object{}

		if len(arguments) > fixed {
			rest = append(rest, arguments[fixed:]...)
		}

		env.Set(function.Parameters[fixed].Value, &object.List{Elements: rest})
	}

	return env, nil
}

// checkArity gives a user-defined function or method the same strict arity
// checking every library call already has (§14 decision 1): a call with the
// wrong number of arguments is an Argument fault naming the call, rather
// than silently dropping extras or leaving a missing parameter undefined. A
// parameter with a default is optional, whichever position it is declared
// in, so the minimum is however many fixed parameters have none; a rest
// parameter (§12) never counts toward the minimum and removes the maximum
// entirely, the same as a library method's own `arityAtLeast`.
func checkArity(function *object.Function, arguments []object.Object, name string, tok token.Token) *object.Error {
	fixed := len(function.Parameters)

	if function.Rest {
		fixed--
	}

	min := fixed - len(function.Defaults)

	if function.Rest {
		return object.ArityAtLeast(name, tok, arguments, min)
	}

	if min == fixed {
		return object.Arity(name, tok, arguments, min)
	}

	return object.ArityRange(name, tok, arguments, min, fixed)
}
