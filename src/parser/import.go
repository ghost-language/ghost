package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) importStatement() ast.ExpressionNode {
	statement := &ast.Import{Token: parser.currentToken}

	parser.readToken()

	if !parser.currentTokenIs(token.STRING) {
		return parser.importFromStatement(statement)
	}

	statement.Path = &ast.String{Token: parser.currentToken, Value: parser.currentToken.Literal.(string)}

	return statement
}

func (parser *Parser) importFromStatement(parent *ast.Import) ast.ExpressionNode {
	statement := &ast.ImportFrom{Token: parent.Token}

	statement.Identifiers = make(map[string]*ast.Identifier)

	if !parser.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	for !parser.currentTokenIs(token.FROM) {
		identifier := &ast.Identifier{Value: parser.currentToken.Lexeme}
		alias := parser.currentToken.Lexeme

		parser.readToken()

		if parser.currentTokenIs(token.AS) {
			parser.readToken()

			alias = parser.currentToken.Lexeme

			parser.readToken()
		}

		statement.Identifiers[alias] = identifier

		if parser.currentTokenIs(token.COMMA) {
			parser.readToken()
		}
	}

	if !parser.currentTokenIs(token.FROM) {
		return nil
	}

	if !parser.expectNextTokenIs(token.STRING) {
		return nil
	}

	statement.Path = &ast.String{Token: parser.currentToken, Value: parser.currentToken.Literal.(string)}

	return statement
}
