package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) superExpression() ast.ExpressionNode {
	return &ast.Super{Token: parser.currentToken}
}
