// Package optimizer rewrites the AST after parsing and before evaluation.
//
// The tree-walking evaluator re-evaluates every node each time control reaches
// it, so an expression built entirely from literals is recomputed on every
// iteration of a loop and every call of a function. Folding those expressions
// into a single literal node at parse time removes that work permanently.
//
// The pass is conservative: it only rewrites a node when the rewritten form is
// guaranteed to produce the same value and the same errors as evaluating the
// original. Anything it does not recognize is left untouched, so failing to
// fold is always safe.
package optimizer

import (
	"ghostlang.org/x/ghost/ast"
)

// globalResolver reports whether a name refers to a library global. The library
// package installs it during initialization. It is a hook rather than a direct
// dependency because library/modules already depends on this package, so
// importing library here would close an import cycle.
//
// When it is nil, identifiers are left unclassified and the evaluator falls
// back to consulting the registries itself.
var globalResolver func(name string) bool

// SetGlobalResolver installs the function used to classify identifiers.
func SetGlobalResolver(resolver func(name string) bool) {
	globalResolver = resolver
}

// Optimize rewrites a parsed program in place and returns it.
func Optimize(program *ast.Program) *ast.Program {
	if program == nil {
		return nil
	}

	for index, statement := range program.Statements {
		program.Statements[index] = optimize(statement)
	}

	return program
}

// optimize walks a node, optimizing its children first so that folding sees
// already-folded operands, then attempts to fold the node itself.
func optimize(node ast.Node) ast.Node {
	switch node := node.(type) {
	case *ast.Program:
		return Optimize(node)

	case *ast.Block:
		if node == nil {
			return node
		}

		for index, statement := range node.Statements {
			node.Statements[index] = optimize(statement)
		}

		return node

	case *ast.Expression:
		node.Expression = optimize(node.Expression)

		return node

	case *ast.Identifier:
		// Classify the name once so the evaluator can skip the library
		// registries for ordinary variables.
		if globalResolver != nil {
			if globalResolver(node.Value) {
				node.LibraryBinding = ast.LibraryBindingGlobal
			} else {
				node.LibraryBinding = ast.LibraryBindingLocal
			}
		}

		return node

	case *ast.Infix:
		node.Left = optimize(node.Left)
		node.Right = optimize(node.Right)

		if folded := foldInfix(node); folded != nil {
			return folded
		}

		return node

	case *ast.Prefix:
		node.Right = optimize(node.Right)

		if folded := foldPrefix(node); folded != nil {
			return folded
		}

		return node

	case *ast.Ternary:
		node.Condition = optimize(node.Condition)
		node.IfTrue = optimize(node.IfTrue)
		node.IfFalse = optimize(node.IfFalse)

		return node

	case *ast.Assign:
		node.Value = optimize(node.Value)

		return node

	case *ast.Compound:
		node.Right = optimize(node.Right)

		return node

	case *ast.Call:
		node.Callee = optimize(node.Callee)
		optimizeExpressions(node.Arguments)

		return node

	case *ast.Method:
		node.Left = optimize(node.Left)
		optimizeExpressions(node.Arguments)

		return node

	case *ast.New:
		node.Class = optimize(node.Class)
		optimizeExpressions(node.Arguments)

		return node

	case *ast.Property:
		node.Left = optimize(node.Left)

		return node

	case *ast.Index:
		node.Left = optimize(node.Left)
		node.Index = optimize(node.Index)

		return node

	case *ast.List:
		optimizeExpressions(node.Elements)

		return node

	case *ast.Map:
		// Map keys are the map's own keys, so a folded key would need the entry
		// reinserted. Values are rewritten in place, which is enough in
		// practice and keeps the pass simple.
		for key := range node.Pairs {
			node.Pairs[key] = optimize(node.Pairs[key])
		}

		return node

	case *ast.If:
		node.Condition = optimize(node.Condition)
		node.Consequence = optimizeBlock(node.Consequence)
		node.Alternative = optimizeBlock(node.Alternative)

		return node

	case *ast.While:
		node.Condition = optimize(node.Condition)
		node.Consequence = optimizeBlock(node.Consequence)

		return node

	case *ast.For:
		if node.Initializer != nil {
			node.Initializer = optimize(node.Initializer)
		}

		if node.Condition != nil {
			node.Condition = optimize(node.Condition)
		}

		if node.Increment != nil {
			node.Increment = optimize(node.Increment)
		}

		node.Block = optimizeBlock(node.Block)

		return node

	case *ast.ForIn:
		node.Iterable = optimize(node.Iterable)
		node.Block = optimizeBlock(node.Block)

		return node

	case *ast.Switch:
		node.Value = optimize(node.Value)

		for _, branch := range node.Cases {
			if branch == nil {
				continue
			}

			optimizeExpressions(branch.Value)

			branch.Body = optimizeBlock(branch.Body)
		}

		return node

	case *ast.Function:
		for name, def := range node.Defaults {
			node.Defaults[name] = optimize(def)
		}

		node.Body = optimizeBlock(node.Body)

		return node

	case *ast.Class:
		node.Body = optimizeBlock(node.Body)

		return node

	case *ast.Trait:
		node.Body = optimizeBlock(node.Body)

		return node

	case *ast.Return:
		if node.Value != nil {
			node.Value = optimize(node.Value)
		}

		return node
	}

	return node
}

func optimizeBlock(block *ast.Block) *ast.Block {
	if block == nil {
		return nil
	}

	for index, statement := range block.Statements {
		block.Statements[index] = optimize(statement)
	}

	return block
}

func optimizeExpressions(expressions []ast.ExpressionNode) {
	for index, expression := range expressions {
		expressions[index] = optimize(expression)
	}
}
