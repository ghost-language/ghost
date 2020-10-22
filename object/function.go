package object

import (
	"bytes"
	"strings"

	"ghostlang.org/x/ghost/ast"
)

type Function struct {
	Parameters []*ast.IdentifierLiteral
	Body       *ast.BlockStatement
	Defaults   map[string]ast.Expression
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString("function")
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String() + "\n")
	out.WriteString("}\n")

	return out.String()
}
