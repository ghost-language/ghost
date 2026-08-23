package evaluator

import (
	"strings"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

// evaluateTemplateString stitches a template literal's text chunks together
// with its interpolated expressions. Each expression is converted to text the
// same way every native function already does when it stringifies a value
// (console.log, print, string.format): by calling the object's own String()
// representation. A template literal needs no explicit toString() call any
// more than logging a value does.
func evaluateTemplateString(node *ast.TemplateString, scope *object.Scope) object.Object {
	var builder strings.Builder

	builder.WriteString(node.Chunks[0])

	for index, expression := range node.Expressions {
		value := Evaluate(expression, scope)

		if isError(value) {
			return value
		}

		if value == nil {
			return object.NewError(fault.Type, node.Token, "cannot interpolate a value that produced nothing into a template literal")
		}

		builder.WriteString(value.String())
		builder.WriteString(node.Chunks[index+1])
	}

	return &object.String{Value: builder.String()}
}
