package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) declaration() ast.StatementNode {
	statement := parser.statement()

	return statement
}
