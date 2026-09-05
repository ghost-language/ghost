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

Ghost supports: classes with inheritance (`extends`), traits (`trait`/`use`), first-class functions, closures, lists, maps, for/for-in/while loops, switch statements, imports, compound operators (`+=`, `++`, etc.), and backtick template literals with `${}` interpolation.

### List operators

Arithmetic operators on lists are elementwise and broadcast (`[1, 2] * 2`, `[[1, 2], [3, 4]] + [10, 20]`). `==` and `!=` compare contents to any depth. Joining two lists is `concat()`, a method, because the operators are arithmetic — do not reintroduce `+` as concatenation. Ordering (`<`, `>`) between lists is deliberately unsupported: neither an elementwise nor a lexicographic reading is obviously right.

### Template literals

Backtick strings interpolate with `${expr}`, JS-style: `` `count: ${count}` ``. Each interpolated value is converted with the same `String()` representation every other native stringification point already uses (`console.log`, `print`, `string.format`), so no type needs an explicit `toString()` call to appear in a template. This is deliberately kept separate from `+`: plain string concatenation between mismatched types (`"count: " + count`) stays a type error, because operators keep one meaning — a template literal is the fluent way to build a mixed-type string, not a loosening of `+`.

The scanner does the work: a template is scanned as alternating chunk and expression tokens (`TEMPLATESTRING`/`TEMPLATESTRINGEND` for text, ordinary tokens for each `${...}`), tracked through a brace-depth stack (`Scanner.templateDepth`) so a `}` from a nested map literal or block doesn't close the interpolation early, and nesting (`` `${ `${x}` }` ``) falls out of the same stack. `parser.templateLiteral()` assembles the chunks and parsed expressions into `ast.TemplateString`; `evaluateTemplateString` stitches the result together at runtime.

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
constructor runs. A method is a *member*, not a name in scope: a method body —
and a field initializer — resolves through the scope the class was declared in,
so a sibling method is reached as `this.method()` and a method never shadows an
import of the same name (§14 decision 12).

## Error handling

Every failure in Ghost is a `*fault.Fault`: a `Kind`, a sentence, a position with
a width, an optional suggested fix, and the calls that were in flight. Nothing
formats `"%d:%d: ..."` into a string by hand — the renderer in `fault/render.go`
turns a fault into the report below, and it is the only code that decides how a
message looks.

```
type error: cannot use `+` between number and string
 --> example.gs:4:15
  |
4 | total = count + label
  |               ^
  |
  = in sum(), called at example.gs:9:1
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

## Working through SPEC.md §11–15 (the 1.0 gap list)

`SPEC.md` §12 (functionality to add) and §13 (defects to fix) are the punch
list for a stable 1.0; §14 records product-direction decisions the list
raises, and §15 ranks the open §13.13–§13.24 findings against each other. Sessions pick this work up incrementally, one item at a time, so the
following are standing rules rather than one-off notes:

- **One callout per session, start to finish.** Implement it, test it, and
  update `SPEC.md` before ending the session — don't leave an item
  half-done or the spec stale relative to the code. Mark the finished
  callout in §12/§13 (e.g. "— done." with a pointer to where it landed and
  where it's tested) rather than deleting it, so the document stays a
  record of what happened, not just what's left. Add newly-implemented
  methods/behavior to the relevant reference table in §8/§9 too — that
  table is what §12/§13 point back to.
- **Match existing naming before inventing new names.** When a gap item
  overlaps a method another type already has (e.g. `string` gaining what
  `list` already calls `contains`/`slice`), reuse that exact name rather
  than a same-meaning synonym — one vocabulary for one operation, the same
  principle §"Error handling" already applies to error wording.
- **A gap item that lists two spellings for one behavior isn't asking for
  both.** Where §12 writes `contains`/`includes` or `charAt`/`at` or
  `slice`/`substring`, that's flagging that *some* name for the operation is
  missing, not requesting synonym methods. Pick one (matching an existing
  convention if there is one) and write down which was picked and why —
  don't add both spellings "to be safe." Genuinely different operations
  written the same way (`indexOf`/`lastIndexOf`, `padStart`/`padEnd`) are
  not this case — implement both.
- **Follow the position-vs-range bounds rule from §13.6.** A method that
  reads *a position* (`charAt`, list/map/string indexing) is lenient —
  return `null`/`""` out of range, never an error. A method that reads *a
  range* (`slice`) validates both ends and raises an `Index` error out of
  range, matching `list.slice()`. Don't mix the two conventions on a new
  method without a documented reason.
- **String/list methods operate on runes, not bytes**, matching
  `string.length()`'s existing `utf8.RuneCountInString`. A new
  position/range-taking method converts to `[]rune` (or converts a byte
  offset back to a rune count, as `indexOf`/`lastIndexOf` do) rather than
  indexing `str.Value` directly — otherwise it silently disagrees with
  `length()` on what a "character" is for any non-ASCII input.
- **Route every new argument check through `object/arguments.go`'s
  helpers** (`Arity`, `ArityRange`, `NumberArgument`, `StringArgument`, ...)
  — never assert a type inline. This is what keeps a new method's error
  messages reading identically to every existing one.
- **Tests live next to the type they cover**, following
  `evaluator/list_methods_test.go` / `evaluator/map_methods_test.go`:
  one `evaluator/<type>_methods_test.go`, using the `evaluate()` /
  `isStringObject`/`isNumberObject`/`isBooleanObject`/`isErrorObject`
  helpers already in `evaluator/evaluator_test.go`, with a table per return
  type and one `Test<Type>MethodErrors` table for argument-error wording
  (assert the *exact* rendered message, position included).
- **Before ending a session on this list**, run `go build ./...` and
  `go test ./...` clean, then leave the next session a clear entry point —
  the next unclaimed item in §12/§13, in the order those sections already
  list them (§12: roughly descending likelihood of being hit by an ordinary
  user; §13: roughly descending damage) unless the user asks for a specific
  item instead. For §13.13–§13.24 — the findings from building Chisel and
  Studio on Ghost — take the highest open row of §15's priority table
  instead; those callouts are numbered in the order they were reported, not
  the order they should be fixed in.
