package ghost

import (
	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/value"
	"ghostlang.org/x/ghost/version"
)

type Ghost struct {
	FatalError bool
	source     string
	file       string
	Scope      *object.Scope
}

var (
	// Version represents the current version.
	Version = version.Version

	// NULL represents a null value.
	NULL = value.NULL

	// TRUE represents a true value.
	TRUE = value.TRUE

	// FALSE represents a false value.
	FALSE = value.FALSE
)

func New() *Ghost {
	scope := &object.Scope{
		Environment: object.NewEnvironment(),
	}

	ghost := &Ghost{
		Scope: scope,
	}

	ghost.registerEvaluator()

	return ghost
}

func (ghost *Ghost) SetDirectory(directory string) {
	ghost.Scope.Environment.SetDirectory(directory)
}

func (ghost *Ghost) GetDirectory() string {
	return ghost.Scope.Environment.GetDirectory()
}

func (ghost *Ghost) SetSource(source string) {
	ghost.source = source
}

func (ghost *Ghost) SetFile(file string) {
	ghost.file = file
}

func (ghost *Ghost) Execute() object.Object {
	scanner := scanner.New(ghost.source, ghost.file)
	parser := parser.New(scanner)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		logParseErrors(parser.Errors())

		return object.NewError(parser.Errors()[0])
	}

	result := evaluator.Evaluate(program, ghost.Scope)

	if object.IsError(result) {
		log.Error(result.(*object.Error).Message)
	}

	return result
}

func RegisterFunction(name string, function object.GoFunction) {
	library.RegisterFunction(name, function)
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	library.RegisterModule(name, methods, properties)
}

// Create a new function called "Call" that will call the passed function with the (optional) passed arguments.
func (ghost *Ghost) Call(function string, args []object.Object) object.Object {
	return ghost.Scope.Environment.Call(function, args, nil)
}

func (ghost *Ghost) registerEvaluator() {
	evaluatorInstance := evaluator.Evaluate

	object.RegisterEvaluator(evaluatorInstance)
	modules.RegisterEvaluator(evaluatorInstance)
}

func logParseErrors(errors []string) {
	for _, message := range errors {
		log.Error(message)
	}
}
