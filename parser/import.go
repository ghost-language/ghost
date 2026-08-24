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

	// `import image, { Spritesheet } from "lumen:image"` — JS-style combined
	// import: a single bare name immediately followed by `, {` binds the
	// whole module under that name (the way `import "lumen:image" as image`
	// does, except the name is chosen positionally here instead of via an
	// `as`), alongside named exports pulled from the braced list — so a
	// script that needs both the module and one of its members doesn't need
	// two `import`s of the same path. `import a, b from "path"` (another
	// bare name after the comma, no brace) is left alone: that already means
	// a second named import, and still does.
	if parser.currentTokenIs(token.IDENTIFIER) && parser.nextTokenIs(token.COMMA) {
		moduleName := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

		parser.readToken() // current: the comma
		parser.readToken() // current: whatever follows it

		if parser.currentTokenIs(token.LEFTBRACE) {
			return parser.combinedImportStatement(parent, moduleName)
		}

		// Not the combined form after all — an ordinary unbraced named-import
		// list that happens to start with two names. The first name and its
		// comma are already consumed, so seed the list with it and let the
		// shared loop parse the rest.
		statement.Identifiers[moduleName.Value] = moduleName
	}

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

// combinedImportStatement parses the rest of `import name, { a, b } from
// "path"` once the leading name and comma have been recognized as that form:
// the braced named list, then `from "path"`. It reuses parent (the ast.Import
// already allocated by importStatement) rather than ast.ImportFrom, because
// this form binds both the whole module — under moduleName, exactly like an
// `as` alias would — and the named exports in Identifiers, which is what
// evaluateImport already knows how to do for the bare `import "path"` form.
func (parser *Parser) combinedImportStatement(parent *ast.Import, moduleName *ast.Identifier) ast.ExpressionNode {
	parent.Alias = moduleName
	parent.Identifiers = make(map[string]*ast.Identifier)

	everything, _, ok := parser.parseImportNames(parent.Token, parent.Identifiers)

	if !ok {
		return nil
	}

	parent.Everything = everything

	// Current is the closing `}`; step onto what should be `from`.
	parser.readToken()

	if !parser.currentTokenIs(token.FROM) {
		parser.report(parser.currentToken, "expected `from` after the closing `}`, found %s", parser.currentToken.Describe()).
			WithHelp("an import reads `import name, { a, b } from \"module\"`")

		return nil
	}

	if !parser.expectNextTokenIs(token.STRING) {
		return nil
	}

	parent.Path = &ast.String{Token: parser.currentToken, Value: parser.currentToken.Literal.(string)}

	return parent
}

// parseImportNames parses the list of names being imported, the shared shape
// behind `import { a, b as c } from "module"`, its unbraced twin `import a, b
// from "module"`, and the braced list in `import name, { a, b } from
// "module"`. The current token on entry is the first token of the list — `{`
// for a braced list, the first name or `*` otherwise.
//
// The names being imported can optionally be wrapped in `{ }`. Braced or not,
// the name list is parsed the same way; only what ends it differs: a braced
// list ends at `}`, an unbraced one at `from`. On return the current token is
// the last one consumed for the list itself — the closing `}` for a braced
// list, or the `from` that stopped an unbraced one — and braced reports which
// kind it was, since callers that need to look past the list (for `from`)
// advance from there themselves.
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
