package evaluator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
)

var searchPaths []string
var imported []string

func evaluateImport(node *ast.Import, scope *object.Scope) object.Object {
	addSearchPath(scope.Environment.GetDirectory())

	filename := findFile(node.Path.Value)

	if filename == "" {
		return object.NewError("%d:%d: runtime error: no file found at '%s.ghost'", node.Token.Line, node.Token.Column, node.Path.Value)
	}

	// Have we imported this file before? If so, we don't need to do anything
	if hasImported(filename) {
		return nil
	}

	evaluateFile(filename, node.Token, scope)

	addImported(filename)

	return nil
}

func evaluateFile(file string, tok token.Token, scope *object.Scope) object.Object {
	source, err := ioutil.ReadFile(file)

	if err != nil {
		return object.NewError("%d:%d: runtime error: %s", tok.Line, tok.Column, err)
	}

	scanner := scanner.New(string(source))
	parser := parser.New(scanner)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		for _, message := range parser.Errors() {
			log.Error(message)
		}

		return nil
	}

	newScope := &object.Scope{Self: scope.Self, Environment: object.NewEnvironment()}

	return Evaluate(program, newScope)
}

func findFile(name string) string {
	basename := fmt.Sprintf("%s.ghost", name)

	for _, path := range searchPaths {
		file := filepath.Join(path, basename)

		if fileExists(file) {
			return file
		}
	}

	return ""
}

func addSearchPath(path string) {
	searchPaths = append(searchPaths, path)
}

func addImported(path string) {
	imported = append(imported, path)
}

func hasImported(path string) bool {
	var result bool = false

	for _, x := range imported {
		if x == path {
			result = true
			break
		}
	}

	return result
}

func fileExists(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}
