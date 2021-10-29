package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) statement() ast.StatementNode {
	return parser.expressionStatement()
}
