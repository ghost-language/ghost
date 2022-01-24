package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) breakStatement() ast.ExpressionNode {
	return &ast.Break{Token: parser.currentToken}
}
