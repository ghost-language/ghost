package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) nullLiteral() ast.ExpressionNode {
	null := &ast.Null{Token: parser.current()}

	return null
}
