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

	// `import "lumen:image", { Spritesheet }` binds the whole module (as
	// above) and pulls named exports out of that same module, in one
	// statement — so getting both `image` and `Spritesheet` out of
	// `"lumen:image"` doesn't need two separate imports of the same path.
	// Only the braced form is accepted here: unlike `import a, b from
	// "path"`, there is no trailing `from` left to mark where the name list
	// ends.
	if parser.nextTokenIs(token.COMMA) {
		parser.readToken()

		if !parser.expectNextTokenIs(token.LEFTBRACE) {
			return nil
		}

		statement.Identifiers = make(map[string]*ast.Identifier)

		everything, _, ok := parser.parseImportNames(statement.Token, statement.Identifiers)

		if !ok {
			return nil
		}

		statement.Everything = everything
	}

	return statement
}

func (parser *Parser) importFromStatement(parent *ast.Import) ast.ExpressionNode {
	statement := &ast.ImportFrom{Token: parent.Token}

	statement.Identifiers = make(map[string]*ast.Identifier)

	everything, braced, ok := parser.parseImportNames(statement.Token, statement.Identifiers)

	if !ok {
		return nil
	}

	statement.Everything = everything

	if braced {
		// Current is the closing `}`; step onto what should be `from`.
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

// parseImportNames parses the list of names being imported, the shared shape
// behind `import { a, b as c } from "module"`, its unbraced twin `import a, b
// from "module"`, and the names in `import "module", { a, b }`. The current
// token on entry is the first token of the list — `{` for a braced list, the
// first name or `*` otherwise.
//
// The names being imported can optionally be wrapped in `{ }`. Braced or not,
// the name list is parsed the same way; only what ends it differs: a braced
// list ends at `}`, an unbraced one at `from`. On return the current token is
// the last one consumed for the list itself — the closing `}` for a braced
// list (deliberately not stepped past, so a caller for whom that `}` ends the
// whole statement can return with the usual "current token is the last token
// of the expression" convention), or the `from` that stopped an unbraced one
// — and braced reports which kind it was, since callers that need to look
// past the list (for `from`) advance from there themselves.
func (parser *Parser) parseImportNames(tok token.Token, identifiers map[string]*ast.Identifier) (everything bool, braced bool, ok bool) {
	braced = parser.currentTokenIs(token.LEFTBRACE)

	if braced {
		parser.readToken()
	}

	if parser.currentTokenIs(token.STAR) {
		everything = true

		parser.readToken()
	} else if !parser.currentTokenIs(token.IDENTIFIER) {
		parser.report(parser.currentToken, "expected a name to import, found %s", parser.currentToken.Describe())

		return false, braced, false
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
				parser.report(tok, "expected `}` to close the names being imported")
			} else {
				parser.report(tok, "expected `from` after the names being imported").
					WithHelp("an import reads `import name from \"module\"`")
			}

			return false, braced, false
		}

		if !parser.currentTokenIs(token.IDENTIFIER) {
			parser.report(parser.currentToken, "expected a name to import, found %s", parser.currentToken.Describe())

			return false, braced, false
		}

		identifier := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
		alias := parser.currentToken.Lexeme

		parser.readToken()

		if parser.currentTokenIs(token.AS) {
			parser.readToken()

			alias = parser.currentToken.Lexeme

			parser.readToken()
		}

		identifiers[alias] = identifier

		if parser.currentTokenIs(token.COMMA) {
			parser.readToken()
		}
	}

	return everything, braced, true
}
