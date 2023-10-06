package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) thisExpression() ast.ExpressionNode {
	return &ast.This{Token: parser.currentToken}
}
