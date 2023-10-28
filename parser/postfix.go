package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) postfixExpression() ast.ExpressionNode {
	return &ast.Postfix{
		Token:    parser.previousToken,
		Operator: parser.currentToken.Lexeme,
	}
}
