package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) expressionStatement() ast.StatementNode {
	expression := parser.expression()

	if expression != nil {
		return &ast.Expression{Expression: expression}
	}

	return expression
}
