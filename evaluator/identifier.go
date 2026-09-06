package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
)

func evaluateIdentifier(node *ast.Identifier, scope *object.Scope) object.Object {
	// Bindings win over library globals, because a global that cannot be
	// shadowed is not a name, it is a reserved word that never announced
	// itself (§13.25). `console` and `type` are ordinary names in an ordinary
	// namespace: a parameter, a loop variable or an import may carry either,
	// and inside its scope that is what the name means.
	if identifier, ok := scope.Environment.Get(node.Value); ok {
		return identifier
	}

	// Nothing is bound, so a global is what is left. The optimizer marks
	// identifiers that cannot name one, letting an undefined name skip two
	// string-keyed map lookups before it reports; an unoptimized AST is left
	// unmarked and still consults the registries.
	if node.LibraryBinding != ast.LibraryBindingLocal {
		if libraryModule, ok := library.GlobalModule(node.Value); ok {
			return libraryModule
		}

		if libraryFunction, ok := library.GlobalFunction(node.Value); ok {
			return libraryFunction
		}
	}

	return undefined(node.Token, node.Value, scope)
}
