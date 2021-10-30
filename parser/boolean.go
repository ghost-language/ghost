package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) boolean() ast.ExpressionNode {
	return &ast.Boolean{Value: parser.check(token.TRUE)}
}
