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

	// `import "math" as m` binds the whole module to `m` instead of the name
	// derived from its path.
	if parser.nextTokenIs(token.AS) {
		parser.readToken()

		if !parser.expectNextTokenIs(token.IDENTIFIER) {
			return nil
		}

		statement.Alias = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	}

	return statement
}

func (parser *Parser) importFromStatement(parent *ast.Import) ast.ExpressionNode {
	statement := &ast.ImportFrom{Token: parent.Token}

	statement.Identifiers = make(map[string]*ast.Identifier)

	// The names being imported can optionally be wrapped in `{ }`, e.g.
	// `import { pi, e } from "math"`. Braced or not, the name list is parsed
	// the same way; only what ends it differs.
	braced := parser.currentTokenIs(token.LEFTBRACE)

	if braced {
		parser.readToken()
	}

	if parser.currentTokenIs(token.STAR) {
		statement.Everything = true

		parser.readToken()
	} else if !parser.currentTokenIs(token.IDENTIFIER) {
		parser.report(parser.currentToken, "expected a name to import, found %s", parser.currentToken.Describe())

		return nil
	}

	stop := token.FROM

	if braced {
		stop = token.RIGHTBRACE
	}

	// Each turn of this loop has to consume a token, and has to stop at the end
	// of the file. An import written the wrong way round — `from "lib" import
	// x` — would otherwise sit here forever looking for a `from` (or a closing
	// `}`) that has already gone past, which is a worse failure than any
	// error: the program never runs and never says why.
	for !parser.currentTokenIs(stop) {
		if parser.isAtEnd() {
			if braced {
				parser.report(statement.Token, "expected `}` to close the names being imported")
			} else {
				parser.report(statement.Token, "expected `from` after the names being imported").
					WithHelp("an import reads `import name from \"module\"`")
			}

			return nil
		}

		if !parser.currentTokenIs(token.IDENTIFIER) {
			parser.report(parser.currentToken, "expected a name to import, found %s", parser.currentToken.Describe())

			return nil
		}

		identifier := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
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

	if braced {
		// Consumes the `}`, landing on what should be `from`.
		parser.readToken()
	}

	if !parser.currentTokenIs(token.FROM) {
		if braced {
			parser.report(parser.currentToken, "expected `from` after the closing `}`, found %s", parser.currentToken.Describe()).
				WithHelp("an import reads `import { name } from \"module\"`")
		}

		return nil
	}

	if !parser.expectNextTokenIs(token.STRING) {
		return nil
	}

	statement.Path = &ast.String{Token: parser.currentToken, Value: parser.currentToken.Literal.(string)}

	return statement
}
