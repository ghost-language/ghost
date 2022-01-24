package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) continueStatement() ast.ExpressionNode {
	return &ast.Continue{Token: parser.currentToken}
}
