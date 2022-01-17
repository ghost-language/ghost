package evaluator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
)

var searchPaths []string
var imported map[string]*object.Scope

func init() {
	imported = make(map[string]*object.Scope)
}

func evaluateImport(node *ast.Import, scope *object.Scope) object.Object {
	addSearchPath(scope.Environment.GetDirectory())

	filename := findFile(node.Path.Value)

	if filename == "" {
		return object.NewError("%d:%d:%s: runtime error: no file found at '%s.ghost'", node.Token.Line, node.Token.Column, node.Token.File, node.Path.Value)
	}

	// Have we imported this file before? If so, we don't need to do anything
	if hasImported(filename) {
		return nil
	}

	addImported(filename, nil)

	moduleScope := evaluateFile(filename, node.Token, scope)

	if isError(moduleScope) {
		return moduleScope
	}

	addImported(filename, moduleScope.(*object.Scope))

	return nil
}

func evaluateImportFrom(node *ast.ImportFrom, scope *object.Scope) object.Object {
	addSearchPath(scope.Environment.GetDirectory())

	filename := findFile(node.Path.Value)

	if filename == "" {
		return object.NewError("%d:%d:%s: runtime error: no file found at '%s.ghost'", node.Token.Line, node.Token.Column, node.Token.File, node.Path.Value)
	}

	// Have we imported this file before? If so, we don't need to do anything
	if hasImported(filename) {
		moduleScope := imported[filename]

		if node.Everything {
			return importEverything(node, scope, moduleScope)
		}

		for alias, identifier := range node.Identifiers {
			value, ok := moduleScope.Environment.Get(identifier.Value)

			if !ok {
				return object.NewError("%d:%d:%s: runtime error: identifier '%s' not found in module '%s.ghost'", node.Token.Line, node.Token.Column, node.Token.File, identifier.Value, node.Path.Value)
			}

			scope.Environment.Set(alias, value)
		}

		return nil
	}

	addImported(filename, nil)

	moduleScope := evaluateFile(filename, node.Token, scope)

	if node.Everything {
		importEverything(node, scope, moduleScope.(*object.Scope))

		addImported(filename, moduleScope.(*object.Scope))

		return nil
	}

	for alias, identifier := range node.Identifiers {
		value, ok := moduleScope.(*object.Scope).Environment.Get(identifier.Value)

		if !ok {
			return object.NewError("%d:%d:%s: runtime error: identifier '%s' not found in module '%s.ghost'", node.Token.Line, node.Token.Column, node.Token.File, identifier.Value, node.Path.Value)
		}

		scope.Environment.Set(alias, value)
	}

	addImported(filename, moduleScope.(*object.Scope))

	return nil
}

func importEverything(node *ast.ImportFrom, scope *object.Scope, moduleScope *object.Scope) object.Object {
	log.Debug("importing everything")

	for alias, value := range moduleScope.Environment.All() {
		scope.Environment.Set(alias, value)
	}

	return nil
}

func evaluateFile(file string, tok token.Token, scope *object.Scope) object.Object {
	source, err := ioutil.ReadFile(file)

	if err != nil {
		return object.NewError("%d:%d:%s: runtime error: %s", tok.Line, tok.Column, tok.File, err)
	}

	directory := scope.Environment.GetDirectory()
	currentFile := strings.Replace(file, directory+"/", "", 1)

	scanner := scanner.New(string(source), currentFile)
	parser := parser.New(scanner)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		for _, message := range parser.Errors() {
			log.Error(message)
		}

		return nil
	}

	newScope := &object.Scope{Self: scope.Self, Environment: object.NewEnvironment()}
	newScope.Environment.SetDirectory(scope.Environment.GetDirectory())

	result := Evaluate(program, newScope)

	if isError(result) {
		return result
	}

	return newScope
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

func addImported(path string, scope *object.Scope) {
	imported[path] = scope
}

func hasImported(path string) bool {
	_, ok := imported[path]

	return ok
}

func fileExists(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}
