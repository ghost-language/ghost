package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIf(node *ast.If, scope *object.Scope) object.Object {
	condition := Evaluate(node.Condition, scope)

	if isError(condition) {
		return condition
	}

	// Each branch runs in a scope of its own, so a name first assigned inside
	// it does not outlive the branch (§13.15). Assignment still walks outward,
	// so writing to a variable declared before the `if` reaches that variable.
	if object.IsTrue(condition) {
		return evaluateBranch(node.Consequence, scope)
	} else if node.Alternative != nil {
		return evaluateBranch(node.Alternative, scope)
	}

	return value.NULL
}

// evaluateBranch runs one arm of an `if` in a scope of its own.
func evaluateBranch(branch ast.StatementNode, scope *object.Scope) object.Object {
	block := enclose(scope)
	result := Evaluate(branch, block)

	release(scope, block)

	return result
}
