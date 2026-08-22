# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
make run        # Build and run Ghost REPL
make build      # Build for all platforms (mac/linux/windows)
make test       # Run all tests with colored output
go test -v ./evaluator/...  # Run tests for a specific package
```

## Architecture Overview

Ghost is a tree-walking interpreter written in Go. The execution pipeline follows this flow:

**Source Code → Scanner → Parser → AST → Evaluator → Object**

### Core Packages

- **scanner/** - Lexical analysis. Transforms source into tokens. Keywords defined in `scanner/scanner.go:21-48`.
- **token/** - Token type definitions and the Token struct containing lexeme, literal, position info.
- **parser/** - Recursive descent parser using Pratt parsing. Each AST node type has its own parsing function (e.g., `parser/function.go`, `parser/if.go`).
- **ast/** - Abstract syntax tree node definitions. Base interfaces in `ast/ast.go`: `Node`, `StatementNode`, `ExpressionNode`.
- **evaluator/** - Tree-walking evaluation. Main entry point is `Evaluate()` in `evaluator/evaluator.go:15`. Each AST node type has a corresponding `evaluate*` function.
- **object/** - Runtime value types (Number, String, Boolean, List, Map, Function, Class, etc.). The `Object` interface (`object/object.go:15-19`) requires `Type()`, `String()`, and `Method()`.
- **ghost/** - Main Ghost struct that orchestrates the pipeline (`ghost/ghost.go`). Entry point for embedding Ghost in Go applications.

### Key Design Patterns

- **Scope**: Wraps Environment and tracks `Self` for method calls (`object/scope.go`).
- **Environment**: Variable storage with parent chain for lexical scoping (`object/environment.go`).
- **Library system**: Native functions and modules registered via `library.RegisterFunction()` and `library.RegisterModule()`. Built-in modules in `library/modules/`. Argument reading and validation goes through the shared helpers in `library/modules/args.go` (`arity`, `numberAt`, `listAt`, `gatherNumbers`, ...) so that argument errors read the same across every module.
- **Math broadcasting**: The math module is split into a scalar layer (`math.go`), reductions (`math_statistics.go`), and arrays and linear algebra (`math_array.go`). Elementwise operations are written against plain numbers and registered with `registerElementwise`, which lifts them to lists and lists of lists through `broadcast()`; reductions are registered with `registerReduction`, which flattens whatever it is given first. Adding a new elementwise method should be a one-line table entry, not a new set of type assertions.

### Object Method System

All object types implement the `Method(method string, args []Object) (Object, bool)` interface. Methods are defined directly on object types (e.g., string methods in `object/string.go`).

## Language Features

Ghost supports: classes with inheritance (`extends`), traits (`trait`/`use`), first-class functions, closures, lists, maps, for/for-in/while loops, switch statements, imports, and compound operators (`+=`, `++`, etc.).

### Class syntax

Classes follow JavaScript/TypeScript conventions. Methods are declared by name
with no `function` keyword, instances are built with `new`, and `super`
resolves members on the superclass of the class that declared the running
method.

```ghost
trait Loud {
    shout() { return this.speak().toUpperCase() }
}

class Animal {
    legs = 4                       // per-instance field

    constructor(name) {
        this.name = name
    }

    speak() { return "..." }
}

class Dog extends Animal {
    use Loud

    constructor(name) {
        super.constructor(name)
    }

    speak() { return super.speak() + " woof" }
}

new Dog("Fido").shout()
```

Field declarations in a class or trait body are initializers, not shared class
state: they are re-evaluated for each instance, ancestors first, before the
constructor runs. A method body's scope is the class environment, so sibling
methods resolve by bare name.

## Version

Update `version/version.go` when releasing. GoReleaser handles binary distribution.
