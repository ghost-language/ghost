package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) expression() ast.ExpressionNode {
	return parser.assign()
}
