package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

func evaluateForIn(node *ast.ForIn, scope *object.Scope) object.Object {
	iterable := Evaluate(node.Iterable, scope)

	if isError(iterable) {
		return iterable
	}

	switch obj := iterable.(type) {
	case *object.List:
		for k, v := range obj.Elements {
			result, done := evaluateForInBody(node, scope, object.NewInt(int64(k)), v)

			if done {
				return result
			}
		}

		return nil
	case *object.Map:
		for _, pair := range obj.OrderedPairs() {
			result, done := evaluateForInBody(node, scope, pair.Key, pair.Value)

			if done {
				return result
			}
		}

		return nil
	}

	return object.NewError(fault.Type, node.Token, "cannot loop over %s", object.TypeName(iterable)).
		WithHelp("`for ... in` walks a list or a map")
}

// evaluateForInBody runs one iteration of the loop, and reports whether that
// iteration ended the loop along with the value to answer with.
//
// The control variables are bound in an environment created for this
// iteration alone, rather than written into the enclosing scope and restored
// afterwards. Two things follow. The variables still do not leak past the
// loop, which is what the save-and-restore was for and what §8.3 promises.
// And a closure made in the body captures that iteration's bindings, so
// `for (name in list) { handlers.push(function () { return name }) }` gives
// each handler its own `name` rather than one shared binding the loop then
// unbinds (§13.14).
func evaluateForInBody(node *ast.ForIn, scope *object.Scope, key object.Object, val object.Object) (object.Object, bool) {
	iteration := enclose(scope)

	iteration.Environment.Set(node.Key.Value, key)
	iteration.Environment.Set(node.Value.Value, val)

	block := Evaluate(node.Block, iteration)

	release(scope, iteration)

	if isTerminator(block) {
		switch result := block.(type) {
		case *object.Error:
			return result, true
		case *object.Return:
			return result, true
		case *object.Continue:
			//
		case *object.Break:
			return nil, true
		}
	}

	return nil, false
}
