package evaluator

import (
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/object"
)

func Register() {
	evaluator := interpreter.Evaluate

	object.SetEvaluator(evaluator)
}
