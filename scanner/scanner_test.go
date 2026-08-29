package scanner

import (
	"testing"

	"ghostlang.org/x/ghost/source"
	"ghostlang.org/x/ghost/token"
)

func TestScanTokens(t *testing.T) {
	test := struct {
		input    string
		expected []struct {
			expectedType   token.Type
			expectedLexeme string
		}
	}{
		`( ) [ ] { } , . - + ; * % ? : > < >= <= ! != = == "hello world" 42 3.14 6.67428e-11 foo foobar hello1 true false class trait use whilefoo こんにちは 世界 += -= *= /= import from as .. index++ index--`,
		[]struct {
			expectedType   token.Type
			expectedLexeme string
		}{
			{token.LEFTPAREN, "("},
			{token.RIGHTPAREN, ")"},
			{token.LEFTBRACKET, "["},
			{token.RIGHTBRACKET, "]"},
			{token.LEFTBRACE, "{"},
			{token.RIGHTBRACE, "}"},
			{token.COMMA, ","},
			{token.DOT, "."},
			{token.MINUS, "-"},
			{token.PLUS, "+"},
			{token.SEMICOLON, ";"},
			{token.STAR, "*"},
			{token.PERCENT, "%"},
			{token.QUESTION, "?"},
			{token.COLON, ":"},
			{token.GREATER, ">"},
			{token.LESS, "<"},
			{token.GREATEREQUAL, ">="},
			{token.LESSEQUAL, "<="},
			{token.BANG, "!"},
			{token.BANGEQUAL, "!="},
			{token.EQUAL, "="},
			{token.EQUALEQUAL, "=="},
			{token.STRING, "hello world"},
			{token.NUMBER, "42"},
			{token.NUMBER, "3.14"},
			{token.NUMBER, "6.67428e-11"},
			{token.IDENTIFIER, "foo"},
			{token.IDENTIFIER, "foobar"},
			{token.IDENTIFIER, "hello1"},
			{token.TRUE, "true"},
			{token.FALSE, "false"},
			{token.CLASS, "class"},
			{token.TRAIT, "trait"},
			{token.USE, "use"},
			{token.IDENTIFIER, "whilefoo"},
			{token.IDENTIFIER, "こんにちは"},
			{token.IDENTIFIER, "世界"},
			{token.PLUSEQUAL, "+="},
			{token.MINUSEQUAL, "-="},
			{token.STAREQUAL, "*="},
			{token.SLASHEQUAL, "/="},
			{token.IMPORT, "import"},
			{token.FROM, "from"},
			{token.AS, "as"},
			{token.DOTDOT, ".."},
			{token.IDENTIFIER, "index"},
			{token.PLUSPLUS, "++"},
			{token.IDENTIFIER, "index"},
			{token.MINUSMINUS, "--"},
			{token.EOF, ""},
		},
	}

	scanner := New(test.input, "test.gs")

	for _, tok := range test.expected {
		token := scanner.ScanToken()

		if tok.expectedType != token.Type {
			t.Fatalf("token type is wrong. expected=%q, got=%q", tok.expectedType, token.Type)
		}

		if tok.expectedLexeme != token.Lexeme {
			t.Fatalf("token lexeme is wrong. expected=%q, got=%q", tok.expectedLexeme, token.Lexeme)
		}
	}
}

// =============================================================================
// Positions and lexical faults

// Every report Ghost prints is anchored on a token, so a token that does not
// know exactly where it starts and how wide it is puts the caret in the wrong
// place. These cover the shapes that used to be measured backwards from the
// scanner's cursor.
func TestTokenPositions(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []struct {
			lexeme string
			line   int
			column int
			length int
		}
	}{
		{
			name:   "operators and names",
			source: `total = count + 1`,
			expected: []struct {
				lexeme string
				line   int
				column int
				length int
			}{
				{"total", 1, 1, 5},
				{"=", 1, 7, 1},
				{"count", 1, 9, 5},
				{"+", 1, 15, 1},
				{"1", 1, 17, 1},
			},
		},
		{
			name:   "numbers start where they are written",
			source: `x = 3.25`,
			expected: []struct {
				lexeme string
				line   int
				column int
				length int
			}{
				{"x", 1, 1, 1},
				{"=", 1, 3, 1},
				{"3.25", 1, 5, 4},
			},
		},
		{
			name:   "a string covers its quotes",
			source: `name = "ghost"`,
			expected: []struct {
				lexeme string
				line   int
				column int
				length int
			}{
				{"name", 1, 1, 4},
				{"=", 1, 6, 1},
				{"ghost", 1, 8, 7},
			},
		},
		{
			name:   "later lines start over at column one",
			source: "a = 1\nbb = 2",
			expected: []struct {
				lexeme string
				line   int
				column int
				length int
			}{
				{"a", 1, 1, 1},
				{"=", 1, 3, 1},
				{"1", 1, 5, 1},
				{"bb", 2, 1, 2},
				{"=", 2, 4, 1},
				{"2", 2, 6, 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := New(test.source, "test.gs")

			for index, expected := range test.expected {
				scanned := scanner.ScanToken()

				if scanned.Lexeme != expected.lexeme {
					t.Fatalf("token %d: got lexeme %q, expected %q", index, scanned.Lexeme, expected.lexeme)
				}

				if scanned.Line != expected.line || scanned.Column != expected.column || scanned.Length != expected.length {
					t.Errorf("token %q: got line=%d column=%d length=%d, expected line=%d column=%d length=%d",
						scanned.Lexeme, scanned.Line, scanned.Column, scanned.Length,
						expected.line, expected.column, expected.length)
				}
			}
		})
	}
}

// A string that runs off the end of the file used to scan silently, leaving the
// rest of the program inside it.
func TestScannerReportsLexicalFaults(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{"unterminated string", `x = "ghost`, "test.gs:1:5: syntax error: unterminated string"},
		{"unterminated comment", "x = 1\n/* forever", "test.gs:2:1: syntax error: unterminated block comment"},
		{"unexpected character", "x = 1\ny = @", "test.gs:2:5: syntax error: unexpected character `@`"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := New(test.source, "test.gs")

			for scanner.ScanToken().Type != token.EOF {
			}

			faults := scanner.Faults()

			if len(faults) != 1 {
				t.Fatalf("got %d faults, expected 1: %v", len(faults), faults)
			}

			if faults[0].String() != test.expected {
				t.Errorf("got=%q, expected=%q", faults[0].String(), test.expected)
			}
		})
	}
}

// A stray character is stepped over rather than swallowed into the name that
// follows it, so the rest of the line still scans.
func TestScanningContinuesPastAStrayCharacter(t *testing.T) {
	scanner := New("@ total", "test.gs")

	scanned := scanner.ScanToken()

	if scanned.Lexeme != "total" {
		t.Errorf("got=%q, expected=%q", scanned.Lexeme, "total")
	}
}

// A string may hold newlines, and the line counter has to follow them or every
// position after one points at the wrong line.
func TestMultiLineStringsAdvanceTheLine(t *testing.T) {
	scanner := New("x = \"one\ntwo\"\ny = 1", "test.gs")

	for index := 0; index < 3; index++ {
		scanner.ScanToken()
	}

	scanned := scanner.ScanToken()

	if scanned.Lexeme != "y" || scanned.Line != 3 {
		t.Errorf("got lexeme=%q line=%d, expected lexeme=%q line=3", scanned.Lexeme, scanned.Line, "y")
	}
}

// A template literal with no interpolation scans as a single closing chunk;
// each `${...}` splits it into a chunk that hands off to ordinary expression
// tokens and a chunk that resumes template text afterward.
func TestTemplateLiteralTokens(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []struct {
			expectedType   token.Type
			expectedLexeme string
		}
	}{
		{
			name:   "no interpolation",
			source: "`hello`",
			expected: []struct {
				expectedType   token.Type
				expectedLexeme string
			}{
				{token.TEMPLATESTRINGEND, "hello"},
				{token.EOF, ""},
			},
		},
		{
			name:   "single interpolation",
			source: "`count: ${1 + 2}`",
			expected: []struct {
				expectedType   token.Type
				expectedLexeme string
			}{
				{token.TEMPLATESTRING, "count: "},
				{token.NUMBER, "1"},
				{token.PLUS, "+"},
				{token.NUMBER, "2"},
				{token.TEMPLATESTRINGEND, ""},
				{token.EOF, ""},
			},
		},
		{
			name:   "text after an interpolation is preserved literally",
			source: "`${a} units`",
			expected: []struct {
				expectedType   token.Type
				expectedLexeme string
			}{
				{token.TEMPLATESTRING, ""},
				{token.IDENTIFIER, "a"},
				{token.TEMPLATESTRINGEND, " units"},
				{token.EOF, ""},
			},
		},
		{
			name:   "a map literal inside an interpolation does not close it early",
			source: "`${ {\"a\": 1}[\"a\"] }`",
			expected: []struct {
				expectedType   token.Type
				expectedLexeme string
			}{
				{token.TEMPLATESTRING, ""},
				{token.LEFTBRACE, "{"},
				{token.STRING, "a"},
				{token.COLON, ":"},
				{token.NUMBER, "1"},
				{token.RIGHTBRACE, "}"},
				{token.LEFTBRACKET, "["},
				{token.STRING, "a"},
				{token.RIGHTBRACKET, "]"},
				{token.TEMPLATESTRINGEND, ""},
				{token.EOF, ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := New(test.source, "test.gs")

			for _, expected := range test.expected {
				scanned := scanner.ScanToken()

				if scanned.Type != expected.expectedType {
					t.Fatalf("token type is wrong. expected=%q, got=%q", expected.expectedType, scanned.Type)
				}

				if scanned.Lexeme != expected.expectedLexeme {
					t.Fatalf("token lexeme is wrong. expected=%q, got=%q", expected.expectedLexeme, scanned.Lexeme)
				}
			}

			if len(scanner.Faults()) != 0 {
				t.Errorf("unexpected faults: %v", scanner.Faults())
			}
		})
	}
}

// An unterminated template literal is reported rather than left to run off the
// end of the file, matching how an unterminated string is handled.
func TestTemplateLiteralUnterminated(t *testing.T) {
	scanner := New("`hello", "test.gs")

	for scanner.ScanToken().Type != token.EOF {
	}

	faults := scanner.Faults()

	if len(faults) != 1 {
		t.Fatalf("got %d faults, expected 1: %v", len(faults), faults)
	}

	expected := "test.gs:1:1: syntax error: unterminated template literal"

	if faults[0].String() != expected {
		t.Errorf("got=%q, expected=%q", faults[0].String(), expected)
	}
}

// The source is filed as scanning starts so that a failure much later — in the
// evaluator, long after parsing finished — can still quote the line.
func TestScanningRegistersTheSource(t *testing.T) {
	source.Reset()

	New("first\nsecond", "registered.gs")

	line, ok := source.Line("registered.gs", 2)

	if !ok || line != "second" {
		t.Errorf("got=(%q, %v), expected=(%q, true)", line, ok, "second")
	}
}
