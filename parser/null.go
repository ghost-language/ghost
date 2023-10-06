package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) nullLiteral() ast.ExpressionNode {
	return &ast.Null{Token: parser.currentToken}
}
