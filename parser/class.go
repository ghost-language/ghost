package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) classStatement() ast.ExpressionNode {
	class := &ast.Class{Token: parser.currentToken}

	parser.readToken()

	class.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

	if parser.nextTokenIs(token.EXTENDS) {
		parser.readToken()

		if !parser.expectNextTokenIs(token.IDENTIFIER) {
			return nil
		}

		class.Super = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	}

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	class.Body = parser.classBody()

	return class
}

// classBody parses the members of a class or trait declaration.
func (parser *Parser) classBody() *ast.Block {
	block := &ast.Block{Token: parser.currentToken}
	block.Statements = []ast.StatementNode{}

	parser.readToken()

	for !parser.currentTokenIs(token.RIGHTBRACE) && !parser.isAtEnd() {
		block.Statements = append(block.Statements, parser.classMember())

		parser.readToken()
	}

	return block
}

// classMember parses a single entry in a class or trait body. Methods use the
// shorthand `name(parameters) { ... }`; everything else — field declarations,
// `use` statements — parses as an ordinary statement.
func (parser *Parser) classMember() ast.StatementNode {
	if parser.currentTokenIs(token.IDENTIFIER) && parser.nextTokenIs(token.LEFTPAREN) {
		return &ast.Expression{Expression: parser.methodDeclaration()}
	}

	return parser.statement()
}

func (parser *Parser) methodDeclaration() ast.ExpressionNode {
	method := &ast.Function{Token: parser.currentToken}
	method.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

	parser.readToken() // step onto the opening parenthesis

	method.Defaults, method.Parameters, method.Rest = parser.functionParameters()

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	method.Body = parser.blockStatement()

	return method
}
