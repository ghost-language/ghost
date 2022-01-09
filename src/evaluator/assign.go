package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateAssign(node *ast.Assign, scope *object.Scope) object.Object {
	value := Evaluate(node.Value, scope)

	if isError(value) {
		return value
	}

	if node.Name != nil {
		switch this := scope.Self.(type) {
		case *object.Class:
			this.Environment.Set(node.Name.Value, value)
		default:
			scope.Environment.Set(node.Name.Value, value)
		}
	}

	return nil
}
