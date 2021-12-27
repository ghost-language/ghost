package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) stringLiteral() ast.ExpressionNode {
	return &ast.String{Token: parser.currentToken, Value: parser.currentToken.Literal.(string)}
}
