package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) stringLiteral() ast.ExpressionNode {
	return &ast.String{Token: parser.peek(), Value: parser.peek().Literal.(string)}
}
