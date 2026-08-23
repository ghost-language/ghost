package ghost

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/color"
	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
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

	// report is where failures are written, and quiet says whether to write
	// them at all. An embedder that wants to handle errors itself turns the
	// reporting off and reads the returned object instead.
	report io.Writer
	quiet  bool
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
		Scope:  scope,
		report: os.Stderr,
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

// SetReportWriter chooses where failures are written. It defaults to standard
// error, which keeps a script's own output clean when it is piped somewhere.
func (ghost *Ghost) SetReportWriter(writer io.Writer) {
	ghost.report = writer
}

// SetQuiet stops Ghost from printing failures. A Go program embedding Ghost
// often wants to present them itself, and it can: the error object comes back
// from Execute either way, carrying everything a report is built from.
func (ghost *Ghost) SetQuiet(quiet bool) {
	ghost.quiet = quiet
}

// Execute scans, parses, and evaluates the source, reporting whatever goes
// wrong along the way.
//
// Every failure leaves by the same door. Syntax errors, runtime errors, and
// bugs in Ghost itself all arrive here as faults, are written out in full, and
// are handed back as an error object. Nothing below this point prints anything
// of its own, and nothing below this point is allowed to reach the caller as a
// Go panic.
func (ghost *Ghost) Execute() object.Object {
	parsed := parser.New(scanner.New(ghost.source, ghost.file))
	program := parsed.Parse()

	if raised := parsed.Errors(); len(raised) != 0 {
		ghost.Report(raised...)

		return object.NewErrorFrom(raised[0])
	}

	result := ghost.evaluate(optimizer.Optimize(program))

	if failed, ok := result.(*object.Error); ok {
		ghost.Report(failed.Fault)
	}

	return result
}

// Report writes faults out in full, one after another, styled if whatever they
// are being written to can show it.
func (ghost *Ghost) Report(raised ...*fault.Fault) {
	if ghost.quiet || ghost.report == nil {
		return
	}

	profile := color.Detect(ghost.report)

	for _, each := range raised {
		fmt.Fprintln(ghost.report, each.Render(profile))
	}
}

// evaluate runs the program, turning a panic anywhere beneath it into an
// ordinary Ghost error.
//
// Nothing should panic — every failure in the evaluator and the library is
// meant to come back as an error object. This is here because "should" is not
// "does": a bug in Ghost would otherwise reach the reader as a Go traceback
// naming files they have never heard of, which tells them nothing about their
// program and nothing about what to do next. Turning it into an internal error
// at least says where in their code it happened, and asks them to report it.
func (ghost *Ghost) evaluate(program *ast.Program) (result object.Object) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = object.NewErrorFrom(internalFault(recovered))
		}
	}()

	return evaluator.Evaluate(program, ghost.Scope)
}

// internalFault describes a panic that escaped the evaluator. The Go stack is
// what a maintainer needs and what a reader has no use for, so it is attached
// only when GHOST_DEBUG asks for it.
func internalFault(recovered interface{}) *fault.Fault {
	raised := fault.New(fault.Internal, "Ghost stopped unexpectedly: %v", recovered)

	if os.Getenv("GHOST_DEBUG") != "" {
		return raised.WithHelp("this is a bug in Ghost; please report it\n\n%s", debug.Stack())
	}

	return raised.WithHelp("this is a bug in Ghost, not in your program; please report it, and set GHOST_DEBUG=1 for the details to include")
}

func RegisterFunction(name string, function object.GoFunction) {
	library.RegisterFunction(name, function)
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	library.RegisterModule(name, methods, properties)
}

// RegisterFunctionForScheme and RegisterModuleForScheme are RegisterFunction/
// RegisterModule for a Go host that wants its own import prefix rather than
// the standard library's `ghost:` — Lumen registering its own `font` module
// so a script can write `import font from "lumen:font"`, for instance.
func RegisterFunctionForScheme(scheme string, name string, function object.GoFunction) {
	library.RegisterFunctionForScheme(scheme, name, function)
}

func RegisterModuleForScheme(scheme string, name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	library.RegisterModuleForScheme(scheme, name, methods, properties)
}

// Call invokes a function defined in the script, with the (optional) passed
// arguments.
func (ghost *Ghost) Call(function string, args []object.Object) object.Object {
	return ghost.Scope.Environment.Call(function, args, nil)
}

func (ghost *Ghost) registerEvaluator() {
	evaluatorInstance := evaluator.Evaluate

	object.RegisterEvaluator(evaluatorInstance)
	modules.RegisterEvaluator(evaluatorInstance)
}
