package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

func evaluateUse(node *ast.Use, scope *object.Scope) object.Object {
	// check that the scope is a class
	class, ok := scope.Self.(*object.Class)

	if !ok {
		return object.NewError(fault.Syntax, node.Token, "`use` can only appear inside a class or a trait")
	}

	var traits []*object.Trait

	for _, trait := range node.Traits {
		if !scope.Environment.Has(trait.Value) {
			return object.NewError(fault.Name, trait.Token, "trait `%s` is not defined", trait.Value)
		}

		identifier, _ := scope.Environment.Get(trait.Value)

		t, ok := identifier.(*object.Trait)

		if !ok {
			return object.NewError(fault.Type, trait.Token, "cannot use `%s`, which is a %s, not a trait", trait.Value, object.TypeName(identifier))
		}

		traits = append(traits, t)
	}

	class.Traits = append(class.Traits, traits...)

	return nil
}
