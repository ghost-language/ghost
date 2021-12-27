package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) booleanLiteral() ast.ExpressionNode {
	return &ast.Boolean{
		Token: parser.currentToken,
		Value: parser.currentTokenTypeIs(token.TRUE),
	}
}
