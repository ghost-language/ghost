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
- **object/** - Runtime value types (Number, String, Boolean, List, Map, Function, Class, etc.). The `Object` interface requires `Type()`, `String()`, and `Method()`.
- **ghost/** - Main Ghost struct that orchestrates the pipeline (`ghost/ghost.go`). Entry point for embedding Ghost in Go applications, and the one place failures are printed.
- **fault/** - Structured diagnostics: what kind of failure, where, how wide, what would fix it, and which calls were in flight. One renderer produces every message Ghost prints.
- **source/** - A registry of scanned source text, so a failure can quote the line it happened on long after parsing finished.
- **color/** - Terminal styling by role rather than by escape code, and the rules for when styling is appropriate at all.

### Key Design Patterns

- **Scope**: Wraps Environment and tracks `Self` for method calls (`object/scope.go`).
- **Environment**: Variable storage with parent chain for lexical scoping (`object/environment.go`).
- **Library system**: Native functions and modules registered via `library.RegisterFunction()` and `library.RegisterModule()`. Built-in modules in `library/modules/`. Argument reading and validation goes through the shared helpers in `library/modules/args.go` (`arity`, `numberAt`, `listAt`, `gatherNumbers`, ...), which wrap the checks in `object/arguments.go` so that argument errors read the same across every module and every method on a value.
- **Broadcasting**: `object.Broadcast()` (`object/broadcast.go`) applies an operation elementwise across operands using numpy's rules — shapes line up from the right, and an axis of length one (or a missing axis) repeats. It lives in `object` because two callers must agree on it exactly: `evaluateListInfix` (`evaluator/list_infix.go`) for `+ - * / %` on lists, and the math module's `broadcast()` for its methods. `a + b` and `math.add(a, b)` are one operation reached two ways — never give either its own rules.
- **Errors**: see "Error handling" below. Nothing builds an error message by hand, and nothing below `ghost.Execute` prints one.
- **Math module**: split into a scalar layer (`math.go`), reductions (`math_statistics.go`), and arrays and linear algebra (`math_array.go`). Elementwise operations are written against plain numbers and registered with `registerElementwise`, which lifts them through `broadcast()`; reductions are registered with `registerReduction`, which flattens whatever it is given first. Adding a new elementwise method should be a one-line table entry, not a new set of type assertions.

### Object Method System

All object types implement the `Method(method string, tok token.Token, args []Object) (Object, bool)` interface. Methods are defined directly on object types (e.g., string methods in `object/string.go`). The token is the method's own name in the source: it is what lets a method report a bad argument at the call that made it, rather than asserting types inline and taking the interpreter down.

## Language Features

Ghost supports: classes with inheritance (`extends`), traits (`trait`/`use`), first-class functions, closures, lists, maps, for/for-in/while loops, switch statements, imports, and compound operators (`+=`, `++`, etc.).

### List operators

Arithmetic operators on lists are elementwise and broadcast (`[1, 2] * 2`, `[[1, 2], [3, 4]] + [10, 20]`). `==` and `!=` compare contents to any depth. Joining two lists is `concat()`, a method, because the operators are arithmetic — do not reintroduce `+` as concatenation. Ordering (`<`, `>`) between lists is deliberately unsupported: neither an elementwise nor a lexicographic reading is obviously right.

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

## Error handling

Every failure in Ghost is a `*fault.Fault`: a `Kind`, a sentence, a position with
a width, an optional suggested fix, and the calls that were in flight. Nothing
formats `"%d:%d: ..."` into a string by hand — the renderer in `fault/render.go`
turns a fault into the report below, and it is the only code that decides how a
message looks.

```
type error: cannot use `+` between number and string
 --> example.ghost:4:15
  |
4 | total = count + label
  |               ^
  |
  = in sum(), called at example.ghost:9:1
  = help: both sides of `+` have to be the same type
```

### Rules

- **Errors never reach Go.** A Ghost program cannot produce a Go panic, a Go
  traceback, or a silent wrong answer. Library methods read their arguments
  through the shared helpers rather than asserting types inline; the evaluator
  counts call depth so a runaway recursion is reported instead of overflowing
  the Go stack; and `ghost.Execute` recovers anything that still escapes and
  reports it as an internal error asking to be filed as a bug.
- **One vocabulary.** `object.Arity`, `object.NumberArgument`, and the rest of
  `object/arguments.go` are the only place an argument problem is worded.
  `library/modules/args.go` wraps them rather than restating them, so a method
  on a string and a method on `math` describe a bad argument identically. Type
  names in messages are the names `type()` answers with.
- **Kinds, not "runtime error".** Pick the `fault.Kind` that says what sort of
  mistake it is: `Syntax`, `Name`, `Type`, `Argument`, `Index`, `Value`,
  `Property`, `Import`, `System`, `Internal`. A failure that fits none of them
  is usually a message that has not been thought through yet.
- **Messages describe, they do not locate.** Write `` "`%s` is not defined" ``,
  not `"%d:%d: runtime error: unknown identifier: %s"`. The position comes from
  the token passed to `object.NewError(kind, tok, ...)`, and the report prints
  it. Messages start lower-case, end without a full stop, and quote code in
  backticks.
- **Help is for what to do next**, not for a second sentence of explanation.
  Where the interpreter is holding the answer — the names in scope, the methods
  a class has — offer it: `evaluator/errors.go` suggests the nearest name for a
  typo.
- **Color is a role, never an escape code.** Ask `color.Detect(writer)` for a
  profile and call `profile.Error(...)`, `profile.Help(...)`, and so on. A
  report rendered plain is the same text minus the paint.

  Detection answers Colored only on positive evidence, because escapes printed
  somewhere that cannot render them are worse than no colour at all. It honours
  `NO_COLOR`, `TERM=dumb`, `FORCE_COLOR`, `CLICOLOR_FORCE`, and `CLICOLOR=0`;
  beyond that the destination has to be a real terminal — an ioctl on Unix, a
  console handle on Windows, not the character-device guess that also matches
  `/dev/null` — and that terminal has to be able to render ANSI, which on Unix
  means a usable `TERM` and on Windows means successfully turning virtual
  terminal processing on. The per-platform half lives in `color/terminal_*.go`.
- **Only the top prints.** The scanner and parser collect faults; the evaluator
  returns error objects; `ghost.Execute` renders them. A `log.Error` inside the
  pipeline is a bug.

### Positions

Tokens carry `Line`, `Column`, and `Length`, all measured from where the lexeme
starts, so a caret underlines the whole of what went wrong. Errors are reported
at the thing that failed rather than at the punctuation next to it: a bad call
points at the callee's name, a bad method at the method's name.

## Version

Update `version/version.go` when releasing. GoReleaser handles binary distribution.
