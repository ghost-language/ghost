package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"regexp"
)

func evaluateExpressions(expressions []ast.ExpressionNode, scope *object.Scope) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, scope)

		// Check if the evaluated value is a string containing a regular expression pattern
		if str, ok := evaluated.(*object.String); ok {
			// If it is, compile the regular expression and store it in the result slice
			regex, err := regexp.Compile(str.Value)
			if err != nil {
				// If there is an error compiling the regular expression,
				// return an error object
				return []object.Object{object.NewError(err.Error())}
			}
			result = append(result, &object.Regexp{Value: regex})
			continue
		}

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

type Regexp struct {
	Value *regexp.Regexp
}

func (r *Regexp) Type() object.ObjectType { return object.RegexpType }
func (r *Regexp) Inspect() string { return r.Value.String() }
