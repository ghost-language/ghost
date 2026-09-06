package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateFunction(node *ast.Function, scope *object.Scope) object.Object {
	// A named function declared in a class or trait body is a method, and a
	// method is a member rather than a lexical binding: it is reached through
	// `this.name()` (§8.8, §14 decision 12). So it closes over the scope the
	// class was *declared* in, not the class's own member table - which is
	// what keeps a method named `math` from hiding an imported `math` from
	// every method in the class (§13.22), and what makes a bare sibling call
	// an honest "not defined" rather than a call with the wrong receiver
	// (§13.17).
	closure := scope

	if node.Name != nil {
		if declaration, ok := scope.Self.(object.FieldDeclarer); ok {
			if result := checkMethodDeclaration(node, declaration); result != nil {
				return result
			}

			closure = declaration.DeclarationScope()
		}
	}

	// The function closes over this scope, so neither it nor anything
	// enclosing it may be reused once this block ends.
	closure.Environment.Capture()

	function := &object.Function{
		Parameters: node.Parameters,
		Defaults:   node.Defaults,
		Body:       node.Body,
		Scope:      closure,
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

// checkMethodDeclaration reports a method colliding with a field of the same
// name in the same body (§13.18). The two would otherwise coexist silently,
// a property read answering with the field and a call answering with the
// method, because neither lookup path knows the other exists.
//
// A method shadowing an imported module is deliberately *not* checked here:
// since a method closes over the class's declaring scope rather than its
// member table, there is no shadowing left to report (§13.22).
func checkMethodDeclaration(node *ast.Function, declaration object.FieldDeclarer) object.Object {
	if declaration.HasField(node.Name.Value) {
		return memberCollisionError(node.Name.Token, node.Name.Value, "field")
	}

	return nil
}

// createFunctionEnvironment binds arguments to a user-defined function's
// parameters, dropping any beyond the last declared one - see checkArity.
// name and tok name and locate the call for an arity error the same way
// every library call already reports one.
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

// checkArity gives a user-defined function or method the same missing-
// argument checking every library call already has (§14 decision 1,
// revised): a call that leaves a required parameter unbound is an Argument
// fault naming the call, rather than leaving that parameter undefined. A
// parameter with a default is optional, whichever position it is declared
// in, so the minimum is however many fixed parameters have none.
//
// There is deliberately no maximum: a caller may pass more arguments than a
// function declares parameters for, and the extras are dropped rather than
// rejected, the same as `object.Function.Evaluate` (used for callbacks
// passed to `list.map`/`filter`/`reduce`/`each`/`sort`) has always allowed.
// This lets a function declare only the parameters its body actually uses -
// a map callback can be `(item) => ...` without also naming the index and
// list every call site provides - instead of forcing every unused trailing
// parameter to be spelled out.
func checkArity(function *object.Function, arguments []object.Object, name string, tok token.Token) *object.Error {
	fixed := len(function.Parameters)

	if function.Rest {
		fixed--
	}

	min := fixed - len(function.Defaults)

	return object.ArityAtLeast(name, tok, arguments, min)
}
