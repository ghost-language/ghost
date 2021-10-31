package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) identifierLiteral() ast.ExpressionNode {
	identifier := &ast.Identifier{Token: parser.current(), Value: parser.current().Lexeme}

	return identifier
}
