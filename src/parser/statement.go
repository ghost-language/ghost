package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) statement() ast.StatementNode {
	statement := parser.assign()

	if statement != nil {
		return statement
	}

	return parser.expressionStatement()
}

func (parser *Parser) expressionStatement() ast.StatementNode {
	statement := &ast.Expression{}
	statement.Expression = parser.parseExpression(LOWEST)

	return statement
}
