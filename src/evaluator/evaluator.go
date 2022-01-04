package evaluator

import (
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
)

func Register() {
	evaluator := interpreter.Evaluate

	object.RegisterEvaluator(evaluator)
	modules.RegisterEvaluator(evaluator)
}
