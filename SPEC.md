# Ghost 1.0 Specification

*What Ghost 1.0 is meant to be: the language and standard library a script can
be written against and keep relying on through the rest of the 1.x line —
plus the concrete work still needed to get there.*

This document specifies the target for a stable 1.0. Most of what follows is
already true of the implementation today — verified by reading the code end
to end, package by package, method by method, rather than assumed from the
docs site (which has already drifted from the implementation in the ways §13
catalogs). Where the implementation does not yet meet the goal specified
here, that gap is named explicitly — in §12 for functionality still to add,
§13 for defects still to fix — rather than folded quietly into the reference
sections as though it already worked. §14 makes the remaining
product-direction calls needed to close both gaps, and §15 ranks what is
left, so a session picking this up has one place to look for what to do
next.

Once the code and this document agree everywhere, cut the tag. Until then, a
disagreement between the two is a bug in one of them, not a reason to guess:
check the implementation, then correct whichever of the two was wrong.

Status at the time of writing: `version/version.go` reports `0.30.0`, the
repository is on a release branch named `1.0`, and `README.md` still says
*"Currently in beta, vetting out the language and seeing how it feels
writing/running. Major changes are still possible at this stage."*

---

## 1. Vision

Ghost is a small, embeddable, dynamically-typed scripting language, implemented
as a tree-walking interpreter in Go. It exists to be the scripting layer
*inside* a larger Go program — a game, a tool, an engine — the way Lua serves
that role in C/C++ projects, but with a syntax and a standard library designed
to feel immediately familiar to someone who already knows JavaScript,
TypeScript, or PHP, and a diagnostic experience modeled on what those
communities' best tooling already does well (see §2 for the specific sources).

1.0 is built around three commitments, in priority order:

1. **Expressive, fluent surface syntax.** Method chains read as sentences.
   Classes look like TypeScript. Template literals build strings the way
   template literals do everywhere else. A user should never have to hold
   Ghost's own quirks in their head on top of the domain they're actually
   trying to script.
2. **A standard library that anticipates what its users reach for**, covers it
   completely and predictably, and never surprises them with a partially
   finished corner. `math` and `date` are the current high-water marks: every
   operation a user is likely to want is present, named consistently, and
   documented in the code as part of a deliberate, audited convention (§7)
   rather than accreted ad hoc. Closing the rest of the standard library to
   that same bar is most of what §12 asks for.
3. **Errors that teach rather than dump.** Every failure is a structured
   `fault.Fault` — a kind, a location, an underlined snippet, a call trace,
   and (where one exists) a suggested fix — rendered by exactly one renderer.
   A Ghost program cannot crash the host process with a raw Go panic or an
   unrecovered stack overflow; every failure, including a bug in Ghost itself,
   comes back as a value the embedding program (or a human at a REPL) can see
   and read.

For 1.0, Ghost is still first and foremost the embeddable scripting layer
described above, not a general-purpose, standalone application language
competing with Python or Node in its own right — it has no concurrency
primitives of its own (concurrency comes from the Go host — see §4 and
§9.11), and 1.0 ships no package manager (§4). That is a scope decision for
this release, not a ceiling on what Ghost is meant to become: the intent is
for Ghost to keep growing into a fully-featured general-purpose language —
standalone scripts and applications, a package system, and whatever else that
implies — over the releases that follow 1.0. Nothing in this document treats
that future as already built; where a 1.0 design decision was made
specifically to leave room for it rather than foreclose it (the `ghost:`
import scheme in §8.9 is the clearest example), that is called out where it
happens, not promised as delivered.

## 2. Inspirations

Ghost borrows deliberately rather than inventing gratuitously, and naming the
sources plainly is more useful than gesturing at "familiar syntax" in the
abstract. The following are load-bearing influences, not incidental
resemblances:

- **JavaScript / TypeScript** — class syntax (`class`, `extends`, `new`,
  `this`, `super`), `camelCase`/`StudlyCase` naming, and template literals
  with `${}` interpolation. Where Ghost's own class syntax reads like JS/TS,
  that is the point, not a coincidence.
- **date-fns** — the `date` module (§9.5) rebuilds date-fns's API shape
  directly: a set of pure functions over an immutable instant rather than a
  mutable `Date` class with methods of its own, plus date-fns's own
  pattern-letter format string (`yyyy-MM-dd`) instead of Go's reference-date
  layout.
- **NumPy** — the broadcasting rules for `+ - * / %` on lists, and for the
  `math` module's arithmetic-as-methods (§9.4), are NumPy's rules exactly:
  shapes align from the right, and an axis of length one repeats.
- **Symfony Console** — `console.write()` versus `console.log()` (§9.3)
  follows the `write()`/`writeln()` distinction Symfony's Console component
  draws between output with and without a trailing newline.
- **Rust's compiler, and modern PHP's error pages** — the fault format (a
  kind, a located and underlined snippet, a call trace, a suggested fix; see
  §8.11) takes its shape from what both already do well: an error a reader
  can act on without reaching for a debugger.
- **Laravel** — not syntax borrowed wholesale, but a discipline: a standard
  library that covers a domain completely rather than piecemeal, a written
  naming convention enforced rather than left as precedent (§7's origin), and
  documentation treated as part of the product rather than an afterthought.
  `math` and `date` are the parts of Ghost that already meet this bar; §12 is
  largely the rest of the standard library catching up to them.
- **PHP, more narrowly, on casing** — considered and explicitly rejected:
  Ghost's methods and properties are `camelCase`, not PHP's `snake_case`,
  because the class syntax already commits to the JS/TS convention and a
  language should not mix the two.

## 3. Goals for 1.0

A 1.0 tag is a promise: the surface syntax and the standard library's existing
names are stable, and a script written against 1.0 keeps working across the
`1.x` line (§7 already states the consequence of breaking that promise — "a
renamed method simply stops existing... made deliberately... and called out in
the version's release notes"). Concretely, 1.0 means:

- **The core language is complete for its stated feature set.** Classes,
  traits, closures, the control-flow forms in §8.6, and the operator set in
  §8.4 do what a reader coming from JS/TS/PHP expects them to do, with no
  silent gaps (see §13 for where that is not yet true).
- **Every standard library module is internally consistent** — one naming
  convention, one style of argument validation, one style of error — and
  externally complete for the domain it claims to own. `math` and `date` meet
  this bar today; `http` stays deliberately minimal by design (§14). `file`
  (renamed and substantially expanded from `io` — §9.8) and `path` (new —
  §9.9) close most of the concrete gap this bar used to name, but have not
  had the same exhaustive audit `math`/`date` went through, so treat them as
  close rather than confirmed complete.
- **No script can crash the host.** True in the general case
  (`ghost.Execute` recovers panics — see §8.11); the one live counterexample
  under concurrent use, the RNG data race, is closed (§13.1 — done).
- **The CLI and its own documentation agree with each other and with the
  code.** `-i` now does (§13.4, done); `-t` still doesn't (§13.10).
- **The naming and design conventions in §7 are actually enforced**, not just
  written down — the property-vs-method inconsistencies §7 now resolves are
  exactly the kind of drift a stated convention exists to prevent.

## 4. Non-Goals

Explicitly out of scope, so they are not mistaken for gaps:

- **Exceptions, `try`/`catch`, or `throw`.** Ghost's error model is values, not
  control flow (§8.11). This is a deliberate, load-bearing design decision
  documented at length in `CLAUDE.md`, not an oversight.
- **Static typing or type annotations.** Ghost is dynamically typed throughout.
- **A package manager or module registry, for 1.0.** `import` resolves either
  a `.gs` file on disk next to the importing file, or a standard library
  name behind the `ghost:` scheme (§8.9); there is no remote fetch, no
  lockfile, no version resolution. This is a scope decision for 1.0, not a
  permanent one (§1) — the `ghost:` scheme exists as a distinct, URL-shaped
  namespace precisely so a future package host can claim its own scheme
  (`pkg:name`, or similar) later without colliding with it or requiring
  `import` itself to change shape.
- **Concurrency primitives in the language itself** (no `async`/`await`, no
  goroutine-equivalent, no channels). Concurrency, where it exists, comes from
  the Go host calling into Ghost from multiple goroutines (`http.handle` is
  the one built-in example) — Ghost code itself is single-threaded per call.
- **Access modifiers, static members, interfaces, or abstract classes.**
  Every class member is public; there is no `static`, `private`, `protected`,
  `interface`, or `abstract` keyword (§8.8). §14 confirms this as a permanent
  stance for the 1.x line, not a placeholder.

## 5. Current Status

| | |
|---|---|
| Version | `0.30.0` (`version/version.go`) |
| Stability | Beta, per `README.md` |
| Release branch | `1.0` |
| Language | Go (`go 1.21.1`, per `go.mod`) |
| External dependencies | `github.com/peterh/liner` (REPL line editing) only |
| Distribution | GitHub Releases via GoReleaser, Homebrew tap, `go install` |
| CI | GitHub Actions: format check, vet, build, test, benchmarks (single-iteration, to catch parse breaks), GoReleaser config validation |
| Test coverage | Present for scanner, parser, evaluator (per-feature files), object equality/broadcast, math, date, JSON, OS — no end-to-end coverage of the CLI flags themselves (which is how §13.4's `-i` gap went unnoticed) |

## 6. Architecture

Ghost is a classic tree-walking interpreter pipeline, plus a constant-folding
pass and a structured-diagnostics layer that both cross-cut every stage:

```
Source text ──▶ Scanner ──▶ Parser ──▶ AST ──▶ Optimizer ──▶ Evaluator ──▶ Object
                  │            │                                 │
                  └────────────┴──────────── fault.Fault ─────────┘
                         (scanner/parser errors)   (runtime errors)
                                       │
                                 source registry
                              (quotes the failing line)
```

| Package | Responsibility |
|---|---|
| `token` | Token type constants and the `Token` struct (lexeme, literal, line/column/length, file). |
| `scanner` | Hand-written lexer. Produces tokens on demand (`ScanToken()`); collects lexical faults but keeps scanning past them. |
| `parser` | Recursive-descent, Pratt-style parser (prefix/infix/postfix function tables keyed by token type, precedence climbing). One file per grammar construct. Collects faults, merges them with the scanner's, and resynchronizes at statement boundaries so one file reports every syntax error in a single pass rather than stopping at the first. |
| `ast` | Plain data: one struct per node kind. `ast.Node`/`StatementNode`/`ExpressionNode`/`AssignmentNode` are empty marker interfaces with no shared fields or methods. |
| `optimizer` | A conservative constant-folding pass (`optimizer.Optimize`) plus one-time classification of every identifier as a library global or an ordinary variable, so the evaluator can skip two map lookups per read of a local. Only rewrites a node when the result is guaranteed identical, including on the error path. |
| `evaluator` | `Evaluate(node, scope)` — one big type switch, delegating to one `evaluate*` function per node kind, in `evaluator/evaluator.go`. |
| `object` | Runtime value types (`Number`, `String`, `Boolean`, `List`, `Map`, `Function`, `Class`, `Instance`, `Trait`, `Super`, `Date`, `Error`, `Scope`, `Environment`, the three `Library*` wrapper types, `NativeClass`). Every value type implements `Method(name, token, args) (Object, bool)`. `new` handles `*Class` directly (building an instance and running its constructor is tree-walking work only the evaluator can do) and falls back to the small `Constructible` interface (`New(scope, token, args...) Object`) for anything else — `*NativeClass` being the one implementation today (§8.9, §10.3). Shared argument-reading/validation helpers live in `object/arguments.go`. |
| `library` | The registry (`library.Functions`, `library.Modules`) that functions and modules install themselves into via `init()`, plus `library.IsGlobal`/`GlobalModule`/`GlobalFunction`, which name the small subset of that registry (`console`, `type`) reachable without an `import` — everything else, built-in or embedder-registered alike, is reached through the `ghost:` import scheme (§8.9, §9.1). |
| `library/functions` | Unqualified library functions. `type` is the one still global (§9.1); a future addition need not be. |
| `library/modules` | The ten built-in modules — see §9. |
| `fault` | The one error model: `Kind`, `Position`, `Frame` (call trace entries), `Fault` itself, and the single renderer (`fault/render.go`) that turns a fault into the boxed, captioned, underlined report every failure in Ghost prints. |
| `source` | A process-wide registry of scanned source text, keyed by filename, so a fault raised long after parsing can still quote the line it happened on. |
| `color` | Decides *whether* output can carry ANSI styling (`color.Detect`, honoring `NO_COLOR`/`TERM=dumb`/`FORCE_COLOR`/`CLICOLOR*`) and exposes styling only by *role* (`Error`, `Help`, `Gutter`, `Caret`, ...), never by raw escape code. |
| `value` | Shared singleton instances: `value.TRUE`, `value.FALSE`, `value.NULL`, `value.BREAK`, `value.CONTINUE`. |
| `ghost` | The embedding facade (`ghost.New()`, `.Execute()`, `.Call()`, `.SetQuiet()`, ...) — the one place a Go host talks to Ghost, and the one place a Go panic anywhere below it is recovered and turned into an internal `fault.Fault` (see §8.11). |
| `log` | A small leveled logger used by the CLI and REPL themselves (not by scripts). |
| `repl` | The interactive shell, built on `github.com/peterh/liner` for line editing and history. |
| `cmd` | The `ghost` binary's `main()` and flag handling. |

### Performance notes worth preserving in any future rewrite

- **Number interning.** `object.Number` stores `int64`/`float64` directly
  (no boxing to `interface{}` internally) and pre-allocates a shared cache of
  small integers (`-128..1024`) so loop counters and indices allocate nothing.
- **Tiered environments.** `object.Environment` stores its first four
  bindings inline in the struct (zero heap allocation for a typical call
  frame), spills to a linear slice up to 16 bindings, and only promotes to a
  real `map[string]Object` beyond that — chosen because a name being looked up
  and the name stored usually share the same backing string from the AST, so
  the scan settles on pointer-equal strings without a byte comparison.
- **Identifier classification at optimize time.** The optimizer marks every
  identifier as a library global or not, once, so the evaluator's hot path
  (`evaluateIdentifier`) skips the two registry lookups entirely for ordinary
  variables.
- **Constant folding.** Literal arithmetic, string concatenation, and boolean
  logic collapse to a single literal AST node before evaluation ever starts,
  so a loop body computing from constants pays for it once, not every
  iteration.

## 7. Naming Conventions

This is the convention every name in Ghost's language and standard library
follows: keywords, library functions and modules, and the methods defined on
each object type. It was written down after an audit found the convention
already in consistent use but nowhere stated — this section is that
statement, so a new addition has a rule to follow instead of just precedent to
imitate.

Ghost's syntax already commits to JavaScript/TypeScript conventions for
classes (`class`, `extends`, `new`), so naming follows that lead rather than
reaching for PHP's `snake_case` (§2). What matters here is the discipline of a
written convention, enforced rather than left as precedent, and comprehensive,
predictable coverage of the standard operations a language's users reach for.

### The rules

**Names are `camelCase`.** Functions, methods, module names, and properties —
`toLowerCase`, `randomInt`, `isPrime`. Classes are `StudlyCase`, matching the
JS/TS convention the class syntax already follows.

**Actions are verbs.** `push`, `pop`, `split`, `trim`, `reshape` — not
`pusher`, `splitting`, `trimmed`. A method name is the imperative form of what
it does.

**Predicates ask a yes/no question and read like one.** Every boolean-
returning library function is named `isX` or `hasX` — `isNaN`, `isEven`,
`isPrime`, `isFinite`, `hasA`. If a name doesn't read as a question, it
shouldn't return a boolean.

**Conversions read left to right, target last.** `toString`, `toNumber`,
`toLowerCase`, `toUpperCase` — the type or form being converted *to* is always
the suffix, so the name reads the same direction as the conversion happens.
`toString()` in particular is universal: every value type answers it,
booleans and null included, so a reader never has to check whether a
particular type happens to support turning itself into a string.

**Modules are lowercase domains, not types.** `math`, `date`, `random`,
`console`, `os` name a *place* an operation belongs, not a kind of value.
Before adding a method to a module, ask which domain the operation is really
in — randomness belongs in `random`, not `math`; a calendar instant belongs in
`date`, not `os`; pausing or ending the program itself belongs in `os`, not
`date`, which is why `os.sleep()` lives beside `os.exit()` rather than beside
`date.now()` — rather than adding a method to whichever module is already
open.

**A name means one thing.** If two modules define a method with the same
name, they'd better mean the exact same operation with the exact same
arguments — the way `a + b` and `math.add(a, b)` are one operation reached two
ways (see `object.Broadcast`). If the arguments or the answer differ, the
names have to differ too, however closely related the operations are.
`math.randomInt()` and `random.random()` are the example this rule is named
for: both draw from the same generator, but one always answers a whole number
and the other a float, so they are not the same name.

**A getter and its setter don't share a name.** A property that can also be
set through a method call needs two different words — one to read it, one to
change it (`random.currentSeed` the property, `random.seed()` the method) —
so that which one a piece of code means is never a guess from context alone.

**A read-only accessor is a property when it costs nothing to answer, and a
method when it doesn't.** A property answers a stored or trivially-derived
value with no meaningful computation behind it — `os.name`, `ghost.version`,
`random.currentSeed`. A method is used when the read does real work or walks
something — `os.args()` builds a fresh list from the process's argv,
`ghost.identifiers()` walks the current scope; both correctly take `()`
because a reader should be able to tell "this is free" from "this does
something" without opening the implementation.

**Errors describe, they don't locate.** This is a naming discipline as much
as an error-message one: `` "`%s` is not defined" ``, not a sentence built
around a position. See §8.11 for the full rules — they're the same spirit
applied to sentences instead of identifiers.

### What counts as a genuine collision

Not every case where two things could theoretically share a word is a
problem. A rename or a consolidation is worth making when:

- Two names mean different things a reader can't tell apart without opening
  both definitions (this is what happened to `math.random`/`random.random`
  before the audit — same word, incompatible argument shapes).
- A name's contract is silently different from what an identical name means
  elsewhere in the language (console's output methods bypassing the writer
  `print()` respects, before the audit fixed it).
- A capability lives in the wrong domain entirely (current-time precision
  living partly in `os`, before it moved to `date`).

It is not a problem when two access points deliberately share an
implementation and agree completely on what they mean — that's the same
pattern as arithmetic operators and the math module's methods, and it doesn't
need fixing, only documenting. Two names checked against this list while
writing this specification, both confirmed as fine as they stand:

- `console.read([prompt])` and `file.read(path)` share a word for operations in
  clearly different domains — interactive input versus a file — with
  different arguments and different failure behavior. A reader is not going
  to mistake one for the other; no rename needed.
- `math.degrees(x)`/`math.radians(x)` are unit conversions of a number to a
  number, the same shape as `sqrt`/`abs`, not a change of value *type* the way
  `toString`/`toNumber` are. The `toX` convention above governs the latter,
  not the former, so these two are correctly named as they are.

### Breaking changes

Ghost has no deprecation mechanism today: a renamed method simply stops
existing under its old name in the next version. A rename is therefore a
breaking change for any `.gs` script written against the old name, made
deliberately rather than as a side effect of a naming pass, and called out in
the version's release notes.

---

## 8. Core Language Reference

### 8.1 Lexical Structure

- **Comments:** `#` and `//` for line comments; `/* ... */` for block comments
  (block comments do not nest).
- **Identifiers:** start with a letter (any Unicode letter) or `_`; continue
  with anything that is not whitespace, a bracket/brace/paren, an operator
  character, or `. , ' " ; :`.
- **Numbers:** integer and floating-point literals, with an optional
  fractional part and an optional scientific-notation suffix (`1e10`,
  `1.5e-3`). No hex, octal, binary, or underscore-grouped literals (`0x1F`,
  `1_000_000` are not supported).
- **Strings:** single- or double-quoted, with backslash escapes for `\n \t \r
  \\ \" \'`; any other escaped character is passed through literally
  (backslash included). Strings may span multiple lines literally.
- **Template literals:** backtick-delimited, JS-style, with `${expr}`
  interpolation (see §8.10). Interpolations nest (`` `${ `${x}` }` ``); the
  scanner tracks brace depth per open interpolation so a `}` belonging to a
  nested map literal or block does not close the interpolation early.
- **Keywords** (30, case-sensitive, all lowercase): `and as break case class
  continue default else extends false for from function if import in new
  null or return super switch this trait true use while`.
  `print` is **not** a keyword — see §13.9.
- **Statement separation:** a trailing `;` is optional and consumed when
  present, after *every* statement kind (§13.12, done) — not just an
  assignment, which used to be the only one that reliably worked;
  otherwise a statement simply ends where its grammar says it ends. There is
  still no significant-newline rule and no automatic-semicolon-insertion
  logic: two statements on separate lines with no `;` between them can still
  parse as one when the second line opens with a token the first's
  expression grammar can continue into (`[`, `(`, `.`, `++`/`--`, a binary
  operator) — `x = 1` then `[10, 20, 30]` on the next line parses as
  `x = 1[10, 20, 30]`, one statement, not two. This is deliberately left as
  is (§13.12) rather than made newline-significant, which would be a much
  larger grammar change with its own trade-offs (multi-line fluent method
  chaining, `foo()\n  .bar()`, relies on a statement being able to continue
  across a line break); ending a line with `;` whenever the next one could
  be read as a continuation avoids it in practice, and now works reliably
  everywhere it didn't before.

### 8.2 Values and Types

`type(value)` and every error message use these exact names:

| Type name | Object | Literal / construction |
|---|---|---|
| `boolean` | `object.Boolean` | `true`, `false` |
| `number` | `object.Number` | integer and float literals; internally distinguishes int64 vs. float64, see below |
| `string` | `object.String` | `"..."`, `'...'`, backtick template literals |
| `list` | `object.List` | `[1, 2, 3]` |
| `map` | `object.Map` | `{key: value, ...}` |
| `null` | `object.Null` | `null` |
| `function` | `object.Function` | `function(...) {...}`, `function name(...) {...}`, a class method |
| `class` | `object.Class` | `class Name {...}` |
| `instance` | `object.Instance` | `new Name(...)` |
| `trait` | `object.Trait` | `trait Name {...}` |
| `date` | `object.Date` | only ever produced by the `date` module |
| `duration` | `object.Duration` | only ever produced by `date.duration()`/`date.durationBetween()` |
| `error` | `object.Error` | any failed operation; also the value a Ghost script cannot construct directly (see §9.12) |
| `super` | `object.Super` | the value of a bare `super` expression |
| `scope` | `object.Scope` | not user-constructible; exists as a `Type()` for internal completeness |
| `library_function`, `library_module`, `library_property` | wrapper types | the value of a bare reference to a module or global function, e.g. `x = math` |
| `break`, `continue`, `return` | internal control-flow markers | never observable as a value inside a running script |

**Numbers are one type with two internal representations.** `object.Number`
holds either an `int64` or a `float64` (never both), tracked by an internal
flag; `type()` reports both as `number`. Integer arithmetic stays exact and
integral; **division (`/`) always promotes to float**, even `4 / 2`, which
answers `2` printed but is a float internally (this matches Python 3's `/`,
not most C-family languages'). Mixed int/float arithmetic promotes to float.
Small integers (`-128..1024`) are interned singletons.

### 8.3 Variables and Scoping

There is no declaration keyword — `x = 5` both declares and assigns, and
which of the two it does depends on whether the name is already in scope.
**Assignment walks the enclosing chain**: it rebinds the nearest existing
binding of the name, wherever that binding lives, and creates a new one — in
the current scope — only when the name is bound nowhere. A function can
therefore update a variable declared outside it, and a block can update one
declared before it, while a genuinely new name stays local to where it first
appears (§14 decision 9, §13.13). Every spelling of assignment follows this
one rule: plain `=`, both destructuring forms, the compound operators, and
`++`/`--`.

Scoping is lexical: closures capture their defining environment, and blocks
(the body of an `if` or `else`, a `while`, a `for`, a `for ... in`, a matched
`switch` case, or a function) each introduce a new enclosed `Environment`
chained to its parent. A name first assigned inside a block does not outlive
it (§13.15) — to keep a value computed in a branch, assign the name before
the branch:

```ghost
result = null
if (ready) { result = compute() } else { result = fallback() }
```

A `for`/`for ... in` loop binds its control variable(s) **once per
iteration**, in that iteration's own scope. Two things follow. The control
variable neither leaks a stray binding into surrounding code nor disturbs an
existing variable of the same name outside the loop. And a closure created in
the loop body captures the value its own iteration ran with, rather than
sharing one binding with every other iteration (§13.14).

**Destructuring assignment** (§12) binds several names from one value in a
single statement: `[a, b] = list` binds positionally (`b` is `null` if
`list` has fewer than two elements, the same leniency `list[i]` itself has
out of range); `{x, y} = map` binds each name from the identically-named map
key (shorthand, matching map literal shorthand keys), and `{x: a} = map`
binds map key `x` to the local name `a` instead. Both are statement-level
only (not usable as an expression, e.g. as a call argument), only bind plain
names (no nesting, and no index/property target — `[a, obj.x] = list` is not
supported), and the right-hand side has to already be a list or a map
respectively, or it is a type error. There is no chained assignment (`a = b
= 5` does not parse as one assignment to both).

### 8.4 Operators

| Category | Operators | Notes |
|---|---|---|
| Arithmetic | `+ - * / %` | On numbers: standard, with the int/float promotion rules above. On lists: elementwise with **NumPy-style broadcasting** — see below. On strings: only `+` (concatenation); `-`/`*`/`/`/`%` on strings are a type error. |
| Comparison | `< <= > >=` | Numbers and strings only (strings compare lexicographically). **Not supported between two lists** — deliberately: neither an elementwise nor a lexicographic reading was judged obviously correct (`CLAUDE.md`). Dates support `< <= > >=` as instant ordering, independent of which time zone either `Date` is attached to (§9.5). |
| Equality | `== !=` | See §8.5 — this is one of the language's most distinctive behaviors. |
| Logical | `and`, `or`, `!` | Word operators, not `&& \|\|` — there is no `&&`/`\|\|` token at all. `!` is the only prefix logical operator. **`and` and `or` short-circuit**: the right operand is evaluated only when the left one leaves the answer open, so `false and x` is `false` and `true or x` is `true` without ever reaching `x` (§13.21, §14 decision 11). Both operands are still booleans — an operand that *is* reached and is not one raises a `Type` fault naming the side at fault — but an operand that is never reached is never type-checked. |
| Unary | `-`, `!` | `-` negates a number only. `!` follows Ghost's truthiness rules (§8.5), not "must be boolean." |
| Range | `a..b` | Inclusive integer range, producing a `list`: `1..5` → `[1, 2, 3, 4, 5]`. Descending (`a > b`) produces an empty list rather than counting down. Not foldable at compile time (would require a shared mutable literal). |
| Ternary | `cond ? a : b` | Standard. |
| Assignment | `=` | Also declares. Valid targets: a bare identifier, an index expression (`list[0] =`, `map["k"] =`), or a property expression (`instance.field =`, `map.key =`). |
| Compound assignment | `+= -= *= /=` | **No `%=`.** Desugars to `target = target OP value`. |
| Increment/decrement | `++ --` | Postfix only (`x++`, not `++x`); operates on a variable, an index, or a property holding a number, mutating the target and evaluating to its value *before* the change — `list[i++]` reads index `i` and then advances it, `while (j++ < n)` tests the old `j`, matching C/JS postfix semantics. |
| Indexing | `a[b]` | Lists (integer index), maps (any hashable key), strings (integer index, returns a one-character string). Out-of-range list/string indices and missing map keys all answer `null` rather than erroring — contrast with `list.slice()`, which errors (see §13.6). |
| Member access | `a.b`, `a.b()` | Property read/assignment vs. method call, disambiguated by whether `(` follows. |

**List broadcasting** (`object.Broadcast`, `object/broadcast.go`) is the same
NumPy-derived rule used by both the `+ - * / %` operators on lists *and* the
`math` module's arithmetic-as-methods (`math.add`, `math.multiply`, ...) —
by design, one implementation, reached two ways:

```
[1, 2, 3] * 2                   # [2, 4, 6]              — a scalar repeats
[1, 2, 3] + [10, 20, 30]        # [11, 22, 33]            — paired elementwise
[[1, 2], [3, 4]] + [10, 20]     # [[11, 22], [13, 24]]    — a row repeats down a matrix
```

Joining two lists end-to-end is a **method**, not `+` — `first.concat(second)`
— because the operators are reserved for arithmetic.

### 8.5 Equality and Truthiness

**Truthiness:** `null` and `false` are falsy; an **empty string is falsy**
(the one type-specific truthiness rule beyond the two obvious ones); every
other value — `0`, `0.0`, an empty list, an empty map — is truthy. This
differs from JS (where `0` and `""` are both falsy) and from Python
(where `0`, `""`, `[]`, and `{}` are all falsy); this is intentional and
should be called out plainly wherever Ghost is introduced to someone arriving
from either background.

`object.IsTrue`/`IsFalse` (`object/boolean.go`) are the one place this rule
is decided — `if`/`while`/`for`/the ternary all call `IsTrue` directly, and
`!` (`evaluator/prefix.go`) answers `toBooleanValue(IsFalse(right))`, rather
than any of them re-deriving the rule locally (§13.11, done). Before that
fix, three independent copies of this same switch existed, and one had
already silently drifted — `!` on an empty string answered `false` instead
of `true` since `evaluator/prefix.go` was first written, caught only by the
audit §13.11 called for, not by any test.

**Equality (`==`/`!=`) is comparison, not coercion**, and behaves differently
depending on the pair of types involved:

| Left / Right | Result |
|---|---|
| Same primitive type (number, string, boolean) | Value equality, as expected. |
| Either side `null` | `true` only if *both* sides are `null`; otherwise `false`. This is the one cross-type comparison that is allowed and does not error. |
| Both `list` | **Deep structural equality**, to any depth (`object.ValuesEqual`/`ListsEqual`). |
| Both `map` | **Deep structural equality**, to any depth, same keys each with an equal value (`object.ValuesEqual`/`MapsEqual` — §13.2, done). |
| Both `instance` | **Identity** — two separate instances with identical fields are not `==`. |
| Both `date` | Instant equality, independent of either operand's attached time zone (§9.5). |
| Both `duration` | **Structural equality** — all six components equal (§9.2). Ordering (`< <= > >=`) stays unsupported, the same reasoning as `list`. |
| Both any other same type (`function`, `class`, `trait`, `scope`, a library wrapper type, ...) | **Identity** — the same rule `instance` gets, now applied uniformly rather than erroring (§13.2, done). |
| Different, non-null types (`5 == "5"`, `[1] == {}`) | **Type error**, not `false`. This is deliberate and covered by an explicit test (`evaluator_test.go`), consistent with `CLAUDE.md`'s "operators keep one meaning" principle. §14 confirms this stays for 1.0 — see that section for what it obligates the documentation to do. |

`object.ValuesEqual` is the one place any of this is decided — the evaluator's
`==`/`!=` operator (`evaluator/infix.go`'s `evaluateEquality`) and `list`'s
`contains()`/`unique()` both call it (except for number/string/boolean/date,
which keep their own dedicated infix evaluator for a reason specific to each:
number promotes int/float the way `+` does, and date compares the instant
rather than every field so two zones can still be equal — both compute
answers `ValuesEqual` would also reach, just without a second trip through
it). A value counted as equal inside a list is equal everywhere else `==` is
asked, including a `map`'s values and a `list` nested inside either.

### 8.6 Control Flow

- **`if (cond) { } else if (cond) { } else { }`** — parentheses around the
  condition and braces around every branch are both mandatory; there is no
  brace-less single-statement form.
- **`while (cond) { }`**.
- **`for (i = 0; i < n; i++) { }`** — C-style three-clause form. The increment
  clause accepts a plain assignment (`i = i + 1`) or, for any target an
  assignment can reach — a variable, an index (`list[0]++`), or a property
  (`obj.count++`) — a compound assignment or a postfix increment/decrement,
  the same as anywhere else an expression is allowed.
- **`for (key, value in iterable) { }`** and **`for (value in iterable) { }`**
  — iterates a `list` (key = integer index) or a `map` (key = the map key, in
  insertion order — §13.5, done). Iterating anything else is a type error
  with the help text *"`for ... in` walks a list or a map."*
- **`switch (value) { case a { } case b, c { } default { } }`** — this is a
  **match-expression**, not a C-style switch: there is no fallthrough, no
  `break` is needed or accepted between cases, and a `case` may list several
  comma-separated values that all run the same block. At most one `default`
  is allowed (a second is a parse error). The subject and every case value
  are compared with `object.ValuesEqual` (§13.2's content-or-identity rule,
  the same one `==` uses) and propagate an error the way every other
  expression does, rather than being compared by string representation and
  silently swallowing a failure (§13.3, done).
- **`break`** and **`continue`** — valid inside `while`, `for`, and
  `for ... in`; propagate through nested blocks via the same
  error/return/break/continue "terminator" mechanism used for early return.
- **`return [expr]`** — valid inside a function or method body; an empty
  `return` yields `null`.

### 8.7 Functions

```ghost
function greet(name, greeting = "Hello") {
    return `${greeting}, ${name}!`
}

add = function(a, b) { return a + b }   // anonymous, assignable

function sum(...numbers) {              // rest parameter
    total = 0
    for (n in numbers) { total = total + n }
    return total
}

sum(...[1, 2, 3])                       // spread at a call site
[0, ...[1, 2], 3]                       // spread in a list literal
```

- Functions are first-class values, close over their defining scope, and may
  be named (hoisted into the enclosing scope/class at definition time) or
  anonymous.
- Parameters support **default values** (`name = "Hello"`), evaluated lazily
  per-call in the function's own scope (so a default can reference an
  earlier parameter or an enclosing binding). Defaults are *not* required to
  trail all non-default parameters — the parser does not enforce an
  ordering.
- The **last** parameter may be a **rest parameter** (`...numbers`, §12): it
  collects every argument from its position onward into a list, always a
  list even when nothing was left to collect (an empty one, never absent). A
  rest parameter cannot have a default — it is always optional on its own —
  and a `...` earlier than the last parameter is a syntax error.
- **Spread** (`...expr`, §12) expands a list's elements in place at a call
  site (`f(...args)`, alongside ordinary arguments: `f(1, ...rest, 2)`) or
  inside a list literal (`[...a, ...b, 1]`). `expr` has to evaluate to a
  list; spreading anything else is a type error. Written anywhere else
  (`x = ...list`), `...` is a syntax error rather than a value — it is only
  meaningful as an argument or a list-literal element.
- User-defined functions and methods get the same **missing-argument
  checking** every library function already has (§14 decision 1, revised): a
  call that leaves a required parameter unbound is an `Argument` fault
  naming the call (`` `foo()` expects at least 2 arguments, got 1 ``), the
  same as a library call — no frame is added, since the call never got the
  chance to start running. A parameter with a default is optional; a rest
  parameter has no upper bound and doesn't count toward the minimum. There
  is no maximum for any function: a call may pass more arguments than a
  function declares parameters for, and the extras are silently dropped, so
  a function only has to name the parameters its body actually uses (a
  `list.map` callback can be `(item) => ...` without also declaring the
  index and list every call passes).
- Recursion is bounded at **4096** call frames (`evaluator/call.go`,
  `maxCallDepth`), reported as an ordinary value error rather than a Go
  stack overflow, and tracked per-`Scope` (not a shared global counter) so
  concurrent goroutines (e.g., separate HTTP requests) do not charge one
  request's recursion against another's budget.

### 8.8 Classes, Traits, and Inheritance

```ghost
trait Loud {
    shout() { return this.speak().toUpperCase() }
}

class Animal {
    legs = 4                       // field: re-evaluated per instance

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

- **`class Name [extends Super] { ... }`** — single inheritance only (no
  multiple `extends`, no interfaces).
- **`trait Name { ... }`** declares a mixin; **`use Trait1, Trait2`** inside a
  class body pulls its members in. A class may `use` more than one trait.
  Member resolution order is: the instance's own fields, then the class's
  own members, then its traits' members, then the superclass chain — see
  `object.LookupMember`.
- **Field declarations** (`legs = 4`) are initializers, not shared class
  state — each is *re-evaluated for every new instance*, ancestors first
  (superclass before subclass, a class's own fields after the traits it
  uses), so two instances never accidentally share one mutable list or map
  through a field default the way a naive "class attribute" implementation
  would let them.
- **`constructor(...)`** is the one specially-named method Ghost recognizes;
  declaring a *field* named `constructor` (`constructor = 5`) is a dedicated
  parse-time error rather than a silently-broken class.
- **`this`** refers to the instance (or, inside a class body outside a
  method, the class itself) currently executing. **`super`** resolves
  members starting at the superclass of the class that *declared* the
  currently-running method — not the receiver's own class — which is what
  keeps a `super` call inside an inherited method from resolving back to
  itself. `super`/`this` outside a class context are name errors.
- A method body's scope is the class environment, so **sibling methods call
  each other by bare name** with no `this.` required — but the bare call
  currently loses the receiver, so it fails as soon as the callee touches
  `this` (§13.17); write `this.method()` until that is fixed (though `this.method()`
  also works).
- **No static members, no access modifiers (`public`/`private`/`protected`),
  no `interface`, no `abstract`, no getters/setters syntax.** Every member is
  a plain, public, instance-level field or method — a permanent stance for
  1.0, confirmed in §14, not a placeholder (§4).
- Instantiation is **`new ClassName(args)`**; `ClassName.new()` is explicitly
  rejected at parse time with a message pointing at the correct syntax
  (leftover PHP/old-JS muscle memory is anticipated and corrected, not just
  rejected).
- `ClassName` can be dotted (`new someMap.ClassName(args)`, `new
  audio.Audio(path)`, §8.9/§10.3) and calls or property reads can chain
  directly onto the result (`new m.Point(1).add(2).toString()`), the same as
  a bare class name already allows — the constructor call is always the
  first call in the chain, however many dots or further calls come after it
  (`splitConstructor` in `parser/new.go`).

### 8.9 Modules and Imports

Two things can appear on the right of an `import`, told apart by scheme:

```ghost
import "helpers"                       // a .gs file, whole-file import for side effects
import add, subtract from "math_ext"   // named imports from a .gs file
import add as plus from "math_ext"     // aliased
import * from "math_ext"               // everything

import "ghost:math"                    // the whole standard library module, bound to `math`
import "ghost:math" as m               // aliased
import { pi, sqrt } from "ghost:math"  // named imports from the standard library
import pi, sqrt from "ghost:math"      // unbraced named imports read the same way

import "lumen:font"                    // an embedding host's own scheme, resolved the same way
import { load } from "lumen:font"      // named imports work identically for any scheme

import { Audio } from "lumen:audio"    // a class exported from a module works the same way too
new Audio("path/to/file.mp3")          // `new` on it works exactly like a Ghost-defined class

import image, { Spritesheet } from "lumen:image"  // JS-style: `image` *and* `Spritesheet`, one statement
image.something()                                  // the whole module, same as `import "lumen:image" as image`
new Spritesheet("sheet.png")                       // the named export, same as `import { Spritesheet } from ...`
```

**File imports** (no scheme — any other string) name a `.gs` file:

- A module is a `.gs` file, looked up **next to the file that imports
  it** (there is no notion of a project root, a search path list beyond
  that, or a package registry — see §4).
- Imports are **process-wide and memoized**: a module is read, parsed, and
  evaluated once no matter how many files import it, guarded by a mutex
  (Ghost code can run on more than one goroutine — e.g. concurrent HTTP
  handlers — and this state is shared across all of them).
- **Circular imports are detected** and reported as a named `Import` fault
  describing the cycle, rather than infinite recursion.
- A module's own parse or evaluation failure surfaces as the importer's
  failure, with a call frame naming the import that pulled it in.
- `import name from "..."` for a name the module does not export suggests
  the nearest name it does, the same typo-correction machinery used
  everywhere else (§8.11).
- The combined form works here too: `import helpers, { greet } from
  "helpers"` binds `helpers` (a `Map` of the module's top-level bindings,
  same as the bare form) and `greet` in one statement.

**Scheme imports** name an entry in a `library.Registry` (§6) with a
`scheme:` prefix instead of a file path — any import path matching
`^[A-Za-z][A-Za-z0-9+.-]+:` (two or more letters before the colon, so a
Windows drive letter is never mistaken for one) is a scheme import rather
than a file path, decided by that pattern alone, with no lookup needed to
tell the two apart. `ghost:` is the scheme the standard library itself
registers under — `import "ghost:math"`, `import { pi } from "ghost:math"`
— and is how the whole standard library reaches a script, `console` and
`type` excepted: they are the only names still reachable without an import
(§9.1), and every other module — `math`, `date`, `random`, `os`, `file`,
`path`, `json`, `http`, `ghost` — has to be imported by name before a script
can use it.

`ghost:` is not special-cased, though — it is simply the one scheme Ghost
itself pre-registers. **A Go program embedding Ghost can claim a scheme of
its own** with `library.RegisterModuleForScheme`/`RegisterFunctionForScheme`
(§10.3) — e.g. a host built as "Lumen" registering `RegisterModuleForScheme
("lumen", "font", methods, properties)` so its own scripts write `import
font from "lumen:font"` rather than borrowing Ghost's own `ghost:`
namespace. The existing unscoped `RegisterModule`/`RegisterFunction` still
target `ghost:` specifically, for backward compatibility with embedding code
written before scheme registration existed; reaching for
`RegisterModuleForScheme` instead is how a host gets a namespace that reads
as its own. Either way there is one registry per scheme and one `scheme:`
import mechanism underneath every scheme, so a module becomes importable the
moment it registers, with no separate mechanism to keep in sync — a plugin
loaded partway through a script (`ghost.extend(...)` followed later by
`import "ghost:itsModule"`, §9.12) works precisely because `import` resolves
at the point it runs, not ahead of time, and the same is true of a host
registering a scheme before running a script versus a plugin registering one
while a script is already running.

- The bare forms — `import "scheme:name"` and `import "scheme:name" as
  alias` — bind the module or function itself. For a module this is the
  exact same value a bare `console` resolves to, so dot access afterward
  (`math.pi`, `math.sqrt()`) works unchanged; for a standalone function
  (there is currently only `type` in the standard library, but an embedder
  can register more under its own scheme) it binds the function, directly
  callable.
- The `from` forms — braced or not, `import { pi } from "ghost:math"` and
  `import pi from "ghost:math"` mean the same thing — pull individual
  methods, properties, and classes out of a module by name, exactly like a
  named import from a `.gs` file. A property (`pi`) is evaluated once,
  immediately, at the import — there is no lazy getter to bind, so
  `import { pi } from "ghost:math"` binds a plain number, the same value
  reading `math.pi` would. A class (`import { Audio } from "lumen:audio"`,
  §10.3) binds the class value itself, unevaluated — there is nothing to
  call at import time, only `new` does that, and it works on an imported
  native class exactly as it does on a Ghost-defined one (§8.8). `import *`
  binds every method, property, and class the module has. The `from` form
  only applies to modules; naming a standalone function this way (there
  being nothing on it to destructure) is a dedicated `Import` fault pointing
  at the bare form instead.
- **The two forms combine, JS-style**: `import name, { a, b } from "path"`
  binds the whole module under `name` — exactly what `import "path" as name`
  does, except the name is chosen positionally here instead of via `as` —
  *and* pulls `a`/`b` out of it by name, in one statement. This is the fix
  for needing both the module itself and one of its members:
  `import image, { Spritesheet } from "lumen:image"` binds `image` (so
  `image.something()` works, if the module has any methods of its own) and
  `Spritesheet` (so `new Spritesheet(...)` does too), where before this took
  two separate `import` lines naming the same path. `import name, { * } from
  "path"` combines the whole module with every export. The named part has to
  be braced here, the same as it does standing alone — a bare
  `import a, b from "path"` (no brace after the comma) keeps its existing,
  unrelated meaning of two named imports, told apart from the combined form
  by that brace alone. Since a standalone function has nothing to pull a
  name out of, naming one this way is rejected the same way
  `import { x } from "scheme:someFunction"` already is.
- A name registered under exactly one scheme, written bare with no import
  (`math.pi` with no import at all) is reported as a `Name` fault with help
  naming the exact import to add, not a generic "did you mean" — see §8.11.
  A name two different schemes both register (unusual, but not prevented)
  lists every import that would resolve it, so the fix names the scheme
  instead of guessing one.
- Misspelling the module name itself (`"ghost:mathh"`), or a name inside the
  `from` form the module does not actually export, gets the same
  nearest-match suggestion as every other unresolved name in Ghost (§8.11).
- Naming a scheme nothing has ever registered under (`"nosuchscheme:thing"`)
  is a distinct `Import` fault from a misspelled name within a real scheme —
  the fix is a different scheme prefix, or waiting until whatever registers
  it has run, not a nearby name within it.

### 8.10 Template Literals

```ghost
count = 3
`there are ${count} item${count == 1 ? "" : "s"}`
```

Backtick strings interpolate `${expr}` (any expression, including nested
template literals). Each interpolated value is converted with that value's
own `String()` — the same representation `print`, `console.log`, and
`string.format` all already use — so no type needs an explicit
`.toString()` to appear cleanly in a template. This is deliberately kept
separate from `+`: `"count: " + count` (mismatched types) is still a type
error; the template literal is the fluent way to build a mixed-type string,
not a loosening of what `+` means.

### 8.11 Error Model

Ghost has **no exceptions**. Every failure — a syntax error, a type
mismatch, a bad argument, a missing file — is a `*fault.Fault` wrapped as an
`object.Error` and returned as an ordinary value; every caller up the chain
(the evaluator, library methods, class methods) checks for one and returns
it unchanged until it reaches the top. There is no `try`/`catch`/`throw` in
the language and no way for a script to swallow an error partway through —
by design (`CLAUDE.md`): "nothing builds an error message by hand," and
"only the top prints."

Every fault carries:

- a **`Kind`** — one of `Syntax, Name, Type, Argument, Index, Value,
  Property, Import, System, Internal` — chosen to say what *sort* of mistake
  it is, never a generic "runtime error";
- a **position** (file, line, column, width) so a report can underline the
  exact offending text;
- an optional **`Help`** — what to do next, not a restated explanation;
- a **call trace**, recorded as the error unwinds through Ghost function and
  method calls (capped at 8 frames, with a "and N more calls" summary
  beyond that);
- and, for a `Name`/`Property`/unresolved-import error, an automatic
  **"did you mean `x`?"** suggestion, computed with a real
  Damerau-Levenshtein edit-distance search over every name actually in
  scope (variables, library globals, class members, or a module's own
  methods, depending on context) — tie-broken by shared-prefix length so
  `pii` suggests `pi` over the equally-distant `phi`.

A single renderer (`fault/render.go`) turns any fault into the boxed report
every failure prints, e.g.:

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

**Nothing below `ghost.Execute` is allowed to reach the caller as a Go
panic.** Runaway recursion is counted and reported as an ordinary value
error before the Go stack actually overflows (§8.7); a bug in Ghost itself
that *does* panic is recovered at the top of `Execute` and reported as an
`Internal` fault asking the reader to file a bug, with the Go stack trace
attached only when `GHOST_DEBUG` is set in the environment. The RNG data
race that used to be the one exception to "nothing crashes the host" under
concurrent use is closed (§13.1 — done).

---

## 9. Standard Library Reference

Every library method validates its own arity and argument types through the
shared helpers in `object/arguments.go` (wrapped, for modules, by
`library/modules/args.go`), so a bad call is always reported as an
`Argument` fault naming the call (`` `math.sqrt()` expects 1 argument, got 0
``), never a Go-level crash.

Every section below except §9.3 (`console`) documents a module that has to be
imported before use — `import "ghost:name"` or `import { a, b } from
"ghost:name"`, per §8.9 — because §9.1 is the complete list of what a script
can use without one.

### 9.1 What's Global

Two names are reachable with no `import` at all — everything else in the
standard library is import-only (§8.9), a deliberate change from earlier
0.x releases where the whole standard library sat in global scope. `console`
and `type` stay global because they are used constantly enough, and small
enough in surface, to earn the same standing print/type-checking primitives
have in other scripting languages — not because "frequently used" is a
standing invitation to add a third.

| Name | Kind | Signature | Behavior |
|---|---|---|---|
| `console` | module | — | Interactive/diagnostic I/O — see §9.3. |
| `type` | function | `type(value)` | Returns the value's type name as a `string` — the exact same name used in every error message (§8.2). Arity: exactly 1. |

There is no bare `print` in 1.0: `console.log(...)` is the one way to write a
line to output. The scanner never produced a `print` keyword token to begin
with — `token.PRINT` was dead code, unreachable from any keyword the scanner
actually recognized, and has been deleted (§13.9, done).

### 9.2 Built-in Methods on Core Types

These are called directly on a value (`"hi".toUpperCase()`), independent of
the module system.

**Bounds-checking convention (§13.6).** A method that reads a *position* —
`list[i]`/`map[k]`/`string[i]` indexing, `charAt()`, `list.get`-shaped
lookups, `pop()`/`shift()`/`first()`/`last()`/`tail()` on an empty list — is
lenient: an out-of-range position answers `null` (or `""` for a string),
never an `Index` error, the same way a missing map key does. A method that
reads a *range* — `slice()` on a list or a string, `fill()` — validates both
ends and raises an `Index` error when either falls outside the receiver,
matching `list.slice()`. The distinction is what the argument names, not the
method's arity: a single position has nowhere else to fall but "absent",
while a range's two ends can straddle the receiver in ways worth catching
before they silently produce a shorter (or, without validation, a garbage)
result. A new list/map/string method follows whichever of the two shapes its
argument is before picking a convention on its own.

**`boolean`** — `toString()`.

**`null`** — `toString()` (→ `"null"`).

**`number`**

| Method | Arity | Behavior |
|---|---|---|
| `round([places])` | 0–1 | Rounds to the nearest integer, or to `places` decimal places. An already-integral `Number` is returned unchanged (not converted through a float round-trip). |
| `floor()` | 0 | Rounds down. Integral input returned unchanged. |
| `ceil()` | 0 | Rounds up. Integral input returned unchanged. |
| `abs()` | 0 | Absolute value; integral input stays an integer. |
| `pow(exponent)` | 1 | Raised to `exponent`. An integer receiver raised to a non-negative integer exponent stays an integer (exact, checked for overflow) so the result can index a list; otherwise falls back to floating point, matching `math.pow()`. |
| `clamp(low, high)` | 2 | The receiver, pulled inside `[low, high]`; answers with one of the three values given rather than a computed one, so clamping whole numbers leaves them whole. A `Value` error if `low` is greater than `high`. |
| `isNaN()` / `isFinite()` / `isInfinite()` | 0 | Only ever true for a float (an integer `Number` can't hold any of the three). |
| `isInteger()` | 0 | True for every non-float `Number`, and for a float only if it is finite with no fractional part. |
| `isEven()` / `isOdd()` | 0 | True only for an integral value (§`isInteger()`) of the matching parity. |
| `isNegative()` / `isPositive()` / `isZero()` | 0 | |
| `toString()` | 0 | |

Each mirrors its `math` counterpart (§9.4) exactly, so `n.pow(2)` and
`math.pow(n, 2)` agree for every input — see §12. `sqrt()`, `lerp()`, and
the rest of `math`'s scalar operations were not given instance-method
counterparts: they have no natural "call this on the receiver" reading the
way `abs`/`pow`/`clamp`/the `isX` predicates do, so they stay reachable
only through `math`.

**`string`**

| Method | Arity | Behavior |
|---|---|---|
| `find(pattern)` | 1 | The receiver is the subject searched; `pattern` (the argument) is compiled as a regular expression — `subject.find(pattern)`, matching JS/PHP/Python (§13.7, done). Returns the first match's full text, or `""` if none. |
| `findAll(pattern)` | 1 | As `find`, but returns a `list` of every match's full text, in order — not, as before §13.7's fix, only the first match's own capture groups. |
| `matches(pattern)` | 1 | As `find`, returning a `boolean`. |
| `format(...args)` | any | `fmt.Sprintf`-style formatting using the receiver as the format string and each argument's `String()`. |
| `endsWith(suffix)` | 1 | |
| `startsWith(prefix)` | 1 | |
| `length()` | 0 | Rune count, not byte count. |
| `replace(from, to)` | 2 | Replaces every occurrence (`strings.ReplaceAll`), not just the first. |
| `split(separator)` | 1 | |
| `toLowerCase()` / `toUpperCase()` | 0 | |
| `toString()` | 0 | Returns the receiver unchanged. |
| `toNumber()` | 0 | Parses as an integer first, then a float; a `Value` error (with help text) if neither parses. |
| `trim()` / `trimStart()` / `trimEnd()` | 0 | Whitespace trim. |
| `contains(substring)` | 1 | Substring search; named to match `list.contains()` rather than adding a second `includes` spelling. |
| `indexOf(substring)` / `lastIndexOf(substring)` | 1 | First/last match position, as a rune index (not byte offset); `-1` if not found. |
| `repeat(n)` | 1 | `n` copies concatenated; a `Value` error if `n` is negative. |
| `padStart(length, pad = " ")` / `padEnd(length, pad = " ")` | 1–2 | Grows to `length` runes by repeating `pad` on the named side, truncated to fit; a receiver already at or past `length`, or an empty `pad`, comes back unchanged. |
| `charAt(index)` | 1 | The single-rune string at a rune position; `""` for a position out of range — a *position* read is lenient, per §13.6 (`at` was not added as a second name for the same thing). |
| `slice(start, end = length())` | 1–2 | New string of the runes from `start` up to `end`; out-of-range bounds raise an `Index` error, matching `list.slice()` (a *range* read validates, per §13.6). `substring` was not added as a second spelling. |
| `reverse()` | 0 | New string, runes reversed. |
| `isEmpty()` | 0 | `length() == 0`. |

Where §12 listed two spellings for one behavior (`contains`/`includes`,
`charAt`/`at`, `slice`/`substring`), only one was implemented — matching
the name `list` already uses where one exists (`contains`, `slice`), and
picking the plainer of the two otherwise (`charAt`, not `at`, since `at`
elsewhere implies negative-index wraparound that nothing in Ghost supports
today). `indexOf`/`lastIndexOf` and `padStart`/`padEnd` are genuinely
different operations, not spelling variants, so both were added.

**`list`**

| Method | Arity | Behavior |
|---|---|---|
| `first()` / `last()` | 0 | `null` on an empty list. |
| `tail()` | 0 | Every element but the first; `null` (not `[]`) on an empty list. |
| `concat(other)` | 1 | New list; joining, not arithmetic — see §8.4. |
| `contains(value)` | 1 | Deep/structural equality (`object.ValuesEqual`). |
| `unique()` | 0 | Drops repeats, preserving first-seen order; same equality rule as `contains`. |
| `each(fn)` | 1 | Calls `fn(element, index)` per element for side effects; returns the list itself (chainable). Short-circuits and propagates on the first error `fn` returns. |
| `map(fn)` | 1 | `fn(element, index)` → new list. |
| `filter(fn)` | 1 | `fn(element, index)` → boolean; keeps truthy results. |
| `reduce(fn[, initial])` | 1–2 | `fn(accumulator, element, index)`. Without an initial value, the first element seeds it; an empty list with no initial value is an `Argument` error. |
| `sort([comparator])` | 0–1 | No comparator: works only for a list that is *entirely* numbers or *entirely* strings (natural order); anything else requires a comparator (`(a, b) → negative/zero/positive number`), or is a dedicated `Argument` error explaining why. Returns a new list (stable sort); does not mutate. |
| `reverse()` | 0 | New list. |
| `push(value)` | 1 | Mutates in place; returns the new length. |
| `pop()` | 0 | Mutates in place; removes and returns the last element, `null` if empty. |
| `shift()` | 0 | Mutates in place; removes and returns the first element, `null` if empty. |
| `slice(start[, end])` | 1–2 | New list; `end` defaults to the list's length. **Out-of-range `start`/`end` is an `Index` error** — contrast with `[]` indexing, which is lenient (§13.6). |
| `join(separator)` | 1 | String-joins each element's `String()`. |
| `length()` | 0 | |
| `toString()` | 0 | |
| `indexOf(value)` | 1 | Position of the first element equal to `value` (`object.ValuesEqual`, same rule as `contains`), or `-1`. Value-based search, pairing with the predicate-based `find`/`findIndex`. |
| `find(fn)` / `findIndex(fn)` | 1 | `fn(element, index)` → truthy; answers the first matching element (`null` if none) or its index (`-1` if none). |
| `some(fn)` / `every(fn)` | 1 | `fn(element, index)` → truthy; short-circuits, mirroring `\|\|`/`&&`. |
| `flatten()` | 0 | New list with every nested list's elements spliced in, recursively. No depth argument — one unambiguous behavior rather than a JS-style default depth. |
| `flatMap(fn)` | 1 | `fn(element, index)` → new list, with each result spliced in one level if it is itself a list, else kept as-is. |
| `chunk(size)` | 1 | New list of new lists of at most `size` elements each, in order; the last chunk holds whatever remains. `size` must be positive (`Value` error otherwise). |
| `fill(value[, start[, end]])` | 1–3 | New list with `value` in place of every element from `start` up to `end` (defaulting to the whole list); out-of-range bounds raise an `Index` error, the same range convention `slice()` uses. Does not mutate, matching `sort()`/`reverse()`/`slice()`. |
| `isEmpty()` | 0 | `length() == 0`. |
| `unshift(value)` | 1 | Mutates in place; front-insert, returns the new length — the front-insert counterpart to `push`. |
| `insertAt(index, value)` / `removeAt(index)` | 2 / 1 | Mutate in place. `insertAt` returns the new length like `push`; an out-of-range `index` clamps to the nearest end rather than erroring. `removeAt` returns the removed element, or `null` for an out-of-range index — the same leniency `pop()`/`shift()` give an empty list. |

There is no generic `splice` — `insertAt`/`removeAt` cover the same ground
as two single-purpose methods, matching the existing `push`/`pop`/`shift`
style rather than one call that both removes and inserts.

**`map`**

| Method | Arity | Behavior |
|---|---|---|
| `get(key[, default])` | 1–2 | `default` if given and the key is absent, else `null`. |
| `has(key)` | 1 | True if the key is present, regardless of its value (distinguishes "absent" from "present and `null`"). |
| `set(key, value)` | 2 | Mutates in place; returns the map itself (chainable). |
| `keys()` / `values()` | 0 | Returns a `list`, in insertion order (§13.5, done). |
| `merge(other)` | 1 | New map, holding this map's pairs in their own order followed by `other`'s; on a key collision, `other`'s value wins but the key keeps this map's position for it (same rule a later assignment to the same key would follow, and the same result a plain object spread gives in JS). |
| `length()` | 0 | |
| `remove(key)` | 1 | Mutates in place; answers the value that was stored under `key`, or `null` if the key was not present — the same leniency `pop()`/`shift()` give an empty list. Named to match `list.removeAt()` rather than adding `delete` as a second spelling (§12). Removing a key drops it from insertion order entirely, leaving no gap. |
| `entries()` | 0 | New `list` of `[key, value]` two-element lists, one per entry, in insertion order — for symmetry with `keys()`/`values()`. |

No `forEach` (use `for ... in`).

**`date`** — `toString()` only (ISO-8601/RFC3339, in the date's own attached
time zone — `Z` for the UTC default, an explicit offset otherwise; §9.5).
Every other date operation is a function in the `date` module (§9.5), not a
method, which keeps `Date` itself immutable and side-effect-free and is a
deliberate design choice modeled on `date-fns` (§2) rather than a mutable
built-in `Date` class — see the doc comment on `object.Date`. Dates support
`< <= > >= == !=` directly as operators (instant comparison, independent of
either operand's attached zone) but no arithmetic operators (`date1 + date2`
is a type error, with help text pointing at the `date` module).

**`duration`** — `years()`, `months()`, `days()`, `hours()`, `minutes()`,
`seconds()` (each arity 0, reading one component back out), and
`toString()` (an ISO 8601 duration — `P1Y2M3DT4H5M6S`, `PT0S` for a zero
duration — the same standard-format instinct `date`'s own `toString()`
follows). Unlike `date`, a `duration` reading its own fields is a method
rather than a module function — see `object.Duration`'s doc comment for
why: a `duration` is a small immutable record, not an opaque instant that
needs a whole module standing between it and its data, so exposing its
fields directly reads more like `map`'s or `list`'s methods than `date`'s
"everything is a free function" stance. Building one, computing one from
two `date`s, and applying one back to a `date` are still module functions
(`date.duration()`/`date.durationBetween()`/`date.addDuration()`/
`date.subDuration()`, §9.5) — those are the operations that construct or
consume a value, not read one already in hand. `==`/`!=` compare a
`duration`'s six components directly (structural equality, like `list` -
§8.5), but ordering (`< <= > >=`) is deliberately unsupported: unlike two
instants, "which of these two spans is longer" has no single answer without
a reference date (a month is a different number of days depending on which
month it started from) — the same reason Temporal requires one for calendar
unit comparisons. No arithmetic operators either (`duration1 + duration2` is
a type error).

### 9.3 `console`

Global — no import needed (§9.1).

Interactive/diagnostic I/O. All output-producing methods write through the
current scope's writer (so redirected output, e.g. inside an `http.handle`
callback, is respected — see §9.11), not directly to the process's stdout.

| Method | Arity | Behavior |
|---|---|---|
| `log(...values)` | any | Space-joined, newline-terminated, no prefix. |
| `info(...values)` | any | As `log`, prefixed `info:` (styled). |
| `warn(...values)` | any | As `log`, prefixed `warning:` (styled). |
| `error(...values)` | any | As `log`, prefixed `error:` (styled). |
| `write(...values)` | any | Like `log`, but **no trailing newline** — the `write`/`writeln` distinction Symfony's Console component draws (§2), and the reason it is named apart from `print`. |
| `newLine()` | 0 | Writes a single newline. |
| `read([prompt])` | 0–1 | Reads one line from stdin (printing `prompt` first, if given); `null` at EOF. |
| `clear()` | 0 | Clears the terminal (`cls`/`clear`, shelling out to the OS rather than using an in-process terminal library). |

No properties currently registered on `console` (`ConsoleProperties` exists
but is empty).

### 9.4 `math`

Import: `import "ghost:math"` or `import { sqrt, pi, ... } from "ghost:math"`
(§8.9).

The most complete module in the standard library, and the model for what
"complete" should mean elsewhere. Organized here exactly as the source
organizes it. Every function below broadcasts across lists and lists of
lists via `object.Broadcast` (§8.4) unless noted otherwise — `math.sqrt`
means the same thing whether given a number, a vector, or a matrix.

**Sign and rounding** — `abs`, `sign`, `floor`, `ceil`, `truncate` (each
arity 1), `round(x[, places])` (arity 1–2, decimal places default 0).

**Powers, roots, logarithms** — `sqrt`, `cbrt`, `square`, `reciprocal`,
`exp`, `exp2`, `expm1`, `log2`, `log10`, `log1p` (arity 1); `pow(base,
exponent)`, `hypot(a, b)` (arity 2, keeps an integer result exact when base
and a non-negative integer exponent allow it); `log(x[, base])` (arity 1–2,
natural log without a base).

**Trigonometry** (radians throughout) — `sin`, `cos`, `tan`, `asin`, `acos`,
`atan`, `sinh`, `cosh`, `tanh`, `asinh`, `acosh`, `atanh` (arity 1);
`atan2(y, x)` (arity 2); `degrees(x)`, `radians(x)` (arity 1, unit
conversion — confirmed correctly named in §7, not a `toX` conversion).

**Special functions** — `gamma`, `logGamma`, `erf`, `erfc` (arity 1).

**Arithmetic as methods** (the same operation the infix operators perform,
exposed as callable/composable functions) — `add`, `subtract`, `multiply`,
`divide`, `mod`, `remainder`, `copySign`, `maximum`, `minimum` (arity 2,
each dividing/modulo variant raising a `Value` error on division by zero).

**Predicates** — `isNaN`, `isFinite`, `isInfinite`, `isInteger`, `isEven`,
`isOdd`, `isNegative`, `isPositive`, `isZero` (arity 1); `isClose(a, b[,
tolerance])` (arity 2–3, default tolerance `1e-9`, combining absolute and
relative tolerance so it holds up near zero and for very large values
alike).

**Bounds and interpolation** — `clamp(x, low, high)`, `lerp(from, to,
amount)`, `smoothstep(low, high, x)` (arity 3); `noise(x[, y])` (arity 1–2,
smoothed value noise, deterministic per input — for terrain/organic
variation, distinct from `random`).

**Randomness sharing `random`'s generator** — `randomInt(n)` /
`randomInt(low, high)` (arity 1–2, inclusive whole number; the one way to
get a whole random number, since `random.random()` always answers a float),
`randomSeed(n)` (arity 1, seeds the same generator `random.seed()` seeds —
this shared state is safe for concurrent use, see §13.1).

**Statistics** (reductions — each accepts a flat argument list, a single
list, or a list of lists, flattened before reducing; arity ≥1) — `sum`,
`product`, `mean`, `median`, `mode` (smallest value wins ties), `variance`,
`sampleVariance` (÷ n−1, requires ≥2 values), `standardDeviation`,
`sampleStandardDeviation`, `min`, `max`, `argmin`, `argmax` (first extreme
wins ties), `cumulativeSum`, `cumulativeProduct` (these two return a list
as long as the input, not a scalar), `gcd`, `lcm`.

**Order statistics** — `percentile(values, p)` (0–100 scale),
`quantile(values, q)` (0–1 scale, same interpolation, arity 2),
`sort(values[, descending])` (arity 1–2), `unique(values)` (arity ≥1).

**Whole-number math** — `factorial(n)` (exact while it fits an `int64`,
falls back to `Γ(n+1)` as a float beyond that), `isPrime(n)`,
`combinations(n, k)`, `permutations(n, k)` (arity 2, each falling back to a
float via the log-gamma identity on overflow rather than wrapping).

**Building arrays** — `arange(stop)` / `arange(start, stop[, step])` (stops
*before* `stop`, keeps whole numbers whole), `linspace(start, stop, count)`
(includes both ends), `zeros(...dims)`, `ones(...dims)`, `full(...dims,
fill)` (1 or 2 dimensions), `identity(n)`.

**Rearranging** — `reshape(list, ...dims)` (one dimension may be `-1` to
infer), `flatten(...)` (any nesting → one flat list), `shape(list)`
(dimensions outermost-first, stopping at the first ragged level),
`transpose(matrix)`.

**Vectors** — `dot(a, b)` (alias: `matmul` — vector·vector → number,
matrix·vector → vector, matrix·matrix → matrix, decided by the shapes
given), `cross(a, b)` (2D or 3D only), `outer(a, b)`, `norm(v[, order])`
(default Euclidean; `order = math.infinity` for the max-norm),
`normalize(v)`, `distance(...)` (accepts either coordinates spread across
the call or two point-vectors), `angle(...)` (two 2D points, `atan2`-based,
full-circle).

**Matrices** — `trace(m)`, `determinant(m)`, `inverse(m)`, `solve(a, b)`
(Gauss-Jordan elimination with partial pivoting under the hood; all three
require a square matrix and report a `Value` error on a singular one rather
than returning `NaN`/`Inf`).

**Constants** (properties, not methods) — `pi`, `tau`, `e`, `phi`, `sqrt2`,
`sqrtPi`, `ln2`, `ln10`, `log2e`, `log10e`, `epsilon`, `smallestNumber`,
`largestNumber`, `infinity`, `nan`, `largestInteger`, `smallestInteger`.

### 9.5 `date`

Import: `import "ghost:date"` or `import { now, of, ... } from "ghost:date"`
(§8.9).

Modeled deliberately on `date-fns` (§2): every function takes a date (or two)
and returns a new value; nothing mutates. A `Date` is an instant plus the
time zone it should be read in, defaulting to UTC — `now()`, `today()`,
`of()`, and `fromUnix()` all build one this way. The instant is what `<`,
`>`, and `==` compare (§8.5) and stays independent of the attached zone, so a
comparison is reproducible everywhere the program runs, the same
reproducibility guarantee a seeded `random` run gets; the zone only governs
what reading a *calendar* position out of that instant answers —
`year()`/`hour()`/`weekday()`, `format()`, `isWeekend()`, `startOfDay()`, and
`toString()` (§14 decision 8, revising the earlier UTC-only design).

Time zones are always explicit and named (IANA identifiers like
`America/New_York`, resolved through the tz database Ghost embeds in its own
binary — see `library/modules/date.go`'s package doc comment), never read
from the host machine's configured zone: there is no `date.local()` or
equivalent. That preserves the original reproducibility goal for the one
case that actually threatened it — a script whose output depended on which
machine ran it — while making the zone-explicit case (`"the same script run
anywhere resolves \`America/New_York\` the same way"`) available, which the
UTC-only design didn't need to give up to get the first guarantee.

**Construction/conversion** — `now()`, `today()` (midnight UTC),
`of(year, month, day[, hour, minute, second])` (month is 1–12; an
out-of-range day/hour/minute/second is a `Value` error rather than silently
rolling into the next period), `ofInZone(year, month, day[, hour, minute,
second], zone)` — `of()` with a required, trailing zone argument: the
components are read as civil time *in that zone*, not in UTC and then
relabeled (`ofInZone(2024, 7, 15, 9, 0, 0, "America/New_York")` is a
different instant than `inTimeZone(of(2024, 7, 15, 9, 0, 0), "America/New_York")`
— the latter keeps the UTC instant and only changes how it is read back),
`parseISO(text)` (RFC3339 or bare `YYYY-MM-DD`; an explicit offset in the
text is preserved rather than normalized away), `fromUnix(seconds)`,
`toUnix(date)`, `toUnixNano(date)`, `format(date, pattern)` (date-fns-style
pattern letters, not Go's reference layout — see below; reads the date's own
zone).

**Time zones** — `inTimeZone(date, zone)` moves a `Date` to a named zone
without changing the instant it names, only what every zone-aware read
answers from that point on; `timeZone(date)` answers the zone's IANA name
(`""` for a `Date` built from a bare numeric offset rather than a named
zone — `parseISO("...-05:00")`, say — since there is no name to report);
`zoneOffset(date)` answers the offset from UTC in seconds, east-positive, at
that specific instant — daylight-saving-aware, so the same named zone can
answer differently for two `Date`s a few months apart. An unrecognized zone
name is a `Value` error naming it.

**Arithmetic** — `addDays`/`subDays`, `addWeeks`/`subWeeks`,
`addMonths`/`subMonths` (clamps to the target month's last day rather than
rolling over — Jan 31 + 1 month = Feb 28/29, not Mar 2/3, and keeps the
result's wall-clock reading in the date's own zone across a daylight-saving
change, the way a calendar app's "same time next month" does),
`addYears`/`subYears`, `addHours`/`subHours`, `addMinutes`/`subMinutes`,
`addSeconds`/`subSeconds` (all arity 2: date, count; these shift by a fixed
real duration regardless of zone, so "add 3 hours" always means 3 real
hours, matching `date-fns`).

**Predicates** — `isSameDay(a, b)`, `isWeekend(date)`, `isLeapYear(date)`
(each reads the date's own attached zone — `isWeekend` on an instant that is
Saturday in UTC but already Sunday in Tokyo answers for Tokyo once moved
there with `inTimeZone`).

**Differences** — `differenceInDays`/`Hours`/`Minutes`/`Seconds(a, b)`,
truncated toward zero (so `differenceInDays(a, b) == -differenceInDays(b,
a)` always, never off-by-one from truncation direction); these compare
instants, so which zone either `Date` is attached to makes no difference to
the result.

**Duration** (§9.2, `object.Duration`, inspired by Temporal's `Duration` and
date-fns's `intervalToDuration` — a calendar-and-clock span kept as separate
components, since "1 month" has no fixed length in days) —
`duration(years, months, days[, hours, minutes, seconds])` builds one
directly (arity 3–6, mirroring `of`'s own shape); every given component has
to point the same direction, all positive or all negative (zero components
don't count toward that), or it's a `Value` error. `durationBetween(a, b)`
computes the calendar breakdown of `a - b` — same sign convention as
`differenceInDays` — as whole months (not split further into years and
months separately, so a month-length ambiguity is never resolved two
different ways), then whole days, then an hours/minutes/seconds remainder;
computed in whichever of `a`/`b` is chronologically earlier's attached zone
— the same zone `addDuration` walks in when reconstructing the later one, so
the two stay exact inverses of each other. This sits alongside, not instead
of, the single-unit `differenceInX` family above — "how many whole days
apart" and "the full calendar breakdown" are different questions, not two
spellings of the same one. `addDuration(date, duration)`/
`subDuration(date, duration)` are the one way to apply a `Duration` back to
a `Date`, walking months, then days, then the clock remainder in that same
order — the order that makes `addDuration(a, durationBetween(b, a))`
reconstruct `b` exactly. Reading a `Duration`'s own fields
(`years()`/`months()`/.../`toString()`) is a method, not a module function —
see §9.2's `duration` entry for why that's the one place this module departs
from "everything beyond `toString()` is a free function."

**Period boundaries** — `startOfDay`, `endOfDay`, `startOfWeek`, `endOfWeek`,
`startOfMonth`, `endOfMonth`, `startOfYear`, `endOfYear` — computed in the
date's own attached zone (midnight in Tokyo for a `Date` moved there, not
midnight UTC). `startOfWeek`/`endOfWeek` treat Sunday as the first day of
the week, matching `weekday()`'s own `0 = Sunday` reading below (and
`date-fns`'s default `weekStartsOn`) rather than the ISO 8601 Monday-first
week.

**Components** — `year`, `month`, `day`, `hour`, `minute`, `second`,
`weekday` (0 = Sunday, matching `date-fns`'s `getDay`) — each reads the
date's own attached zone, so the same instant can answer a different day,
hour, or weekday depending on which zone it was last moved to.

**Format pattern letters** (a run of the same letter is one token; anything
else copies through literally): `y` (year), `M` (month, 1/01/Jan/January by
run length), `d` (day, padded at 2+), `E` (weekday, abbreviated below 4
letters), `H`/`h` (24h/12h hour), `m` (minute), `s` (second), `a` (AM/PM).

### 9.6 `random`

Import: `import "ghost:random"` or `import { random, seed } from "ghost:random"`
(§8.9).

| Method/Property | Arity | Behavior |
|---|---|---|
| `random()` / `random(max)` / `random(min, max)` | 0–2 | Uniform float; no args → `[0, 1)`; one arg → `[0, max)`; two args → `[min, max)`. |
| `seed([n])` | 0–1 | Reseeds the shared generator; no argument reseeds from the current time. |
| `currentSeed` *(property)* | — | The seed currently driving the generator. |

`SeedRandom` is also exported at the Go level so an embedding host can fix
reproducibility before a script ever runs. This generator is shared with
`math.randomInt`/`math.randomSeed`; both modules read and write it through
the same mutex-guarded accessors, so it is safe for concurrent use (§13.1 —
done).

### 9.7 `os`

Import: `import "ghost:os"` or `import { sleep, exit } from "ghost:os"`
(§8.9).

| Method/Property | Arity | Behavior |
|---|---|---|
| `args()` | 0 | The process's own command-line arguments (excluding the binary name). |
| `exit(code[, message])` | 1–2 | Prints `message` if given, then terminates the process immediately (`os.Exit` — no deferred cleanup runs). Arguments are validated *before* anything is printed, so a miscall does not partially execute. |
| `sleep(milliseconds)` | 1 | Blocks the calling goroutine; negative durations are a `Value` error. |
| `name` *(property)* | — | The OS name (`runtime.GOOS`: `"linux"`, `"darwin"`, `"windows"`, ...). |

`args()` and `name` correctly follow the property-vs-method rule in §7: one
does real work (builds a list from argv), the other answers a stored value.

### 9.8 `file`

Import: `import "ghost:file"` or `import { read, write } from "ghost:file"`
(§8.9). Renamed from `io` — the old name was a placeholder that never named
what the module actually does (read and write files), and nothing else was
ever added under it to justify the more generic name.

| Method | Arity | Behavior |
|---|---|---|
| `read(path)` | 1 | Reads a whole file as a string. |
| `write(path, content)` | 2 | Overwrites an **existing** file's contents, keeping its permissions; errors (with help text) if the file does not already exist — pointed deliberately at `file.append` for creating one. |
| `append(path, content)` | 2 | Creates the file if needed (mode `0644`); appends `content` plus a trailing newline. |
| `exists(path)` | 1 | `boolean` — whether anything (file or directory) is at `path`. |
| `isDirectory(path)` | 1 | `boolean` — a `System` error, not `false`, if `path` does not exist at all (distinct from `exists`, which answers the existence question itself). |
| `size(path)` | 1 | The file's size in bytes, as a `number`. |
| `delete(path)` | 1 | Removes a file or an **empty** directory (the same reach as `os.Remove` in Go) — a non-empty directory is left alone, a `System` error, rather than wiped. |
| `mkdir(path)` | 1 | Creates a directory, and any missing parent directories along the way (`mkdir -p`), mode `0755`. |
| `list(path)` | 1 | The entry names (not full paths) directly inside a directory, as a `list` of strings — not recursive. |
| `copy(source, destination)` | 2 | Copies a file's contents and permissions to `destination`, which is created or overwritten. |
| `move(source, destination)` | 2 | Renames/moves a file or directory. |

All paths are resolved **relative to the running script's own directory**,
not the process's current working directory (an absolute path is used as
given) — so a script behaves the same regardless of where `ghost` was
invoked from.

Closes the gap the old `io` module left (§12 used to list `exists`,
`delete`/`remove`, `mkdir`, `listDir`/`readDir`, `copy`, `move`/`rename` as
missing — all now present here). Still no streaming read/write and no
recursive `list`/`delete`; a script working with very large files or deep
directory trees has to shell out or wait for a future addition.

### 9.9 `path`

Import: `import "ghost:path"` or `import { join, basename } from "ghost:path"`
(§8.9). Pure string manipulation — nothing here touches the filesystem or
knows where the running script lives, which is the actual line between this
module and `file` (§7: a module names a domain, and *building* a path is a
different domain from *reading* one, even though a script usually does both
together).

| Method | Arity | Behavior |
|---|---|---|
| `join(...parts)` | ≥1 | Joins path segments with the OS's own separator, cleaning the result (`..`/`.` resolved, redundant separators collapsed) the way `filepath.Join` does. |
| `basename(path)` | 1 | The last path element. |
| `dirname(path)` | 1 | Everything but the last path element. |
| `extname(path)` | 1 | The extension, dot included (`".gs"`), or `""` if there is none. |
| `isAbsolute(path)` | 1 | `boolean`. |

### 9.10 `json`

Import: `import "ghost:json"` or `import { decode, encode } from "ghost:json"`
(§8.9).

| Method | Arity | Behavior |
|---|---|---|
| `decode(text)` | 1 | Parses JSON into a `list` or `map`. A syntactically valid JSON document whose top level is a bare scalar (number/string/boolean/null) is rejected with help text suggesting the caller wrap it in `[ ]`. |
| `encode(value)` | 1 | Serializes a `list` or `map` (recursively) to a JSON string. Map keys are converted via their `String()`; any other top-level argument type is an `Argument` error. |

No pretty-printing option (indentation), no streaming, and — because a
`Date` has no JSON representation of its own — encoding a `Date` inside a
structure falls through `object.ObjectToAnyValue`'s unhandled-type case to
Go's `nil`/`null` silently (worth confirming as intended, since it is easy
to lose a date value with no error raised).

### 9.11 `http`

Import: `import "ghost:http"` (§8.9).

| Method | Arity | Behavior |
|---|---|---|
| `handle(path, callback)` | 2 | Registers `callback(request)` for `path` via Go's `net/http.HandleFunc`. `request` is a `map` with `method`, `host`, `contentLength`, `protocol`, `protocolMajor`, `protocolMinor`, `body` (the raw request body as a string). **The callback's own output writer is redirected to the HTTP response** — so a handler produces its response body by calling `console.log`/etc. *inside* the callback, not by returning a value or calling a dedicated "respond" function. A panic inside a handler goroutine is recovered locally and answered with a 500, rather than taking the whole server down. |
| `listen(port[, ready])` | 1–2 | Starts the server (blocking) on `port` (1–65535); `ready`, if given, is called once binding is confirmed set up (before the blocking call, in practice — see note below). Handles `SIGINT` for graceful shutdown (30s timeout). A port that cannot be bound is a `System` error, not a silent hang. |

This module is, and stays, intentionally minimal for 1.0 (§14): a script
cannot set a response status code, set a header, read query-string
parameters, read cookies, or parse a route parameter (`/:id`) —
everything beyond "handle a path, write a body" is out of scope for this
release.

### 9.12 `ghost`

Import: `import "ghost:ghost"` (§8.9) — the module name and the scheme
prefix are spelled the same by coincidence of domain, not because the
prefix is special-cased for it; `import { abort } from "ghost:ghost"` works
the same as any other named import.

Reflective/meta operations over the running Ghost instance itself.

| Method/Property | Arity | Behavior |
|---|---|---|
| `abort(reason)` | 1 | `reason` a `string` raises it as a `Value` error (the one place a *script* raises an error deliberately, rather than an operation failing on its own); `reason` `null` is a silent no-op; anything else is an `Argument` error. |
| `execute(code)` | 1 | Parses and evaluates `code` as a fresh Ghost program in the *current* scope — dynamic `eval`, in effect. Its own syntax errors are reported as this call's failure. |
| `extend(pluginPath)` | 1 | Loads a Go plugin (`.so`) via Go's `plugin` package and calls its exported `Register()` function, which is expected to call back into `library.RegisterFunction`/`RegisterModule`. **Linux/macOS only — the Go `plugin` package has no Windows support**, a real portability gap given Ghost ships Windows binaries (§13.8). |
| `identifiers()` | 0 | Every name currently bound in scope, as a `list` of strings (introspection/debugging aid). |
| `version` *(property)* | — | The running Ghost version string. |

---

## 10. CLI and Embedding

### 10.1 CLI

```
ghost [flags] [file]
```

| Flag | Behavior | Status |
|---|---|---|
| *(none)* | Starts the REPL. | Works. |
| `file` | Executes the file, then exits with status 1 if it failed. | Works. |
| `-h` | Prints usage and exits. | Works. |
| `-v` | Prints the binary name and version. | Works. |
| `-t` | Prints how long execution took. | Works, and documented in `helpCommand()`'s own printed help (§13.10, done). |
| `-i` | Runs a file, then drops into a REPL with the script's environment intact — same `Scope`, so every binding the script made is there to inspect, even one made before a mid-script failure. Composes with `-t` (timing prints before the prompt starts); with no file argument it is a no-op, same as plain `ghost` alone. | Works (§13.4, done). |

### 10.2 REPL

Line editing and history via `github.com/peterh/liner`; history persists to
`$TMPDIR/.ghost_history`. Each line is executed as its own program against a
single, persistent `ghost.Ghost` instance, so state accumulates across
lines the way a REPL is expected to work. `Ctrl+C` aborts the current line;
`Ctrl+D` ends the session; a failed line is reported in full and the session
continues.

### 10.3 Embedding API (`ghost` package)

The primary way a Go program hosts Ghost:

```go
instance := ghost.New()
instance.SetDirectory(dir)      // for import/file path resolution
instance.SetFile(name)          // for error reporting
instance.SetSource(code)
instance.SetQuiet(true)         // suppress ghost's own error printing
instance.SetReportWriter(w)     // where errors print, if not quiet (default stderr)

result := instance.Execute()    // returns object.Object; check object.IsError(result)
value := instance.Call("fnName", []object.Object{...})  // call into the script afterward

ghost.RegisterFunction("myFn", myGoFunc)
ghost.RegisterModule("myModule", methods, properties)

// Or, to claim a namespace of the embedder's own instead of borrowing
// ghost:, e.g. for a host built as "Lumen":
ghost.RegisterFunctionForScheme("lumen", "myFn", myGoFunc)
ghost.RegisterModuleForScheme("lumen", "myModule", methods, properties)

// A class whose instances are built and driven entirely by Go code — a
// stateful host resource (an audio handle, say) the host wants scripts to
// `new` and call methods on, rather than expose as a bag of functions:
ghost.RegisterClassForScheme("lumen", "audio", "Audio", audioConstructor)
// -> import { Audio } from "lumen:audio"
// -> new Audio("path/to/file.mp3")
```

`RegisterFunction`/`RegisterModule`/`RegisterFunctionForScheme`/
`RegisterModuleForScheme` are process-global (they write into shared
registries — `library.Functions`/`library.Modules` for the first two, a
`library.Registry` per scheme for the latter two). §14 decides this is not a
1.0 requirement to change: a host needing isolated configurations runs
separate processes for now.

A module or function registered this way is not automatically global to
scripts — nothing is, `console`/`type` aside (§9.1). It becomes reachable
the same way every built-in module is: `import "ghost:myModule"` or
`import { myFn } from "ghost:myFn"` for the unscoped calls, `import
"lumen:myModule"` for a call made `ForScheme("lumen", ...)` (§8.9).
`RegisterFunction`/`RegisterModule` are the original, unscoped calls and
keep targeting `ghost:` specifically, for embedding code written before
scheme registration existed; `RegisterFunctionForScheme`/
`RegisterModuleForScheme` are how a host claims a scheme that reads as its
own rather than Ghost's. Both pairs write into the same kind of registry and
resolve through the same `import` mechanism (§8.9), so there is no second
convention to learn or keep in sync for a host that wants its own
namespace — only which of the two calls to make.

**`RegisterClass`/`RegisterClassForScheme`** register a class the same way —
a member of a module, alongside whatever methods and properties it may also
have — except its instances are built and driven entirely by Go rather than
by evaluating a Ghost class body. `object.NativeClass` is what gets
registered (`Name` and a `Constructor` with the exact signature `object.
GoFunction` already has); `new` accepts one through a small `object.
Constructible` interface (`New(scope, tok, args...) Object`) — a second,
generic path alongside the dedicated one `*object.Class` already had, not a
parallel class system, since `Constructible` is the only thing `new` needs
to know about a class it didn't build itself. `Constructor` returns whatever
value a `new Audio(...)` should
produce — any `object.Object` at all, with its own `Method()` implementation
deciding what calling something on an instance does, entirely on the host's
side; Ghost does not prescribe an instance shape beyond "is an `object.
Object`". Reading a property or calling a method on the class value itself,
rather than an instance (`Audio.path`, `Audio.path()`), is refused with the
same wording a Ghost-defined class gives for the same mistake — a script
using it has no way to tell, and no reason to need to tell, which kind of
class `Audio` actually is.

This is a Go-level extension point only: there is still exactly one way to
declare a class in `.gs` source (§8.8), and nothing here adds a second.
A bundled Ghost-source "prelude" — classes written in `.gs`, evaluated
once at startup, calling down into native functions for the primitive work —
was considered as an alternative and deliberately not built: it would have
split embedding into two conventions (a Go-level one for functions/modules,
a Ghost-source one for anything class-shaped) for no benefit an embedder
actually needs today. `NativeClass` keeps embedding to the one, Go-level
convention throughout. It is the right tool specifically when the class
wraps a genuinely native operation with no meaningful Ghost-level logic
around it (opening a file handle, decoding audio) — tree-walking a method
that only ever calls straight back into Go buys nothing. A class that is mostly orchestration over calls that are themselves already
native — a hypothetical future `file.File` wrapping `file.read`/`file.write`
(§9.8), say — would be the case for reconsidering a Ghost-source prelude, if
one ever turns up in the standard library itself; nothing here forecloses
that, it just isn't needed yet.

### 10.4 Extending Ghost from Ghost itself

`ghost.extend(path)` (§9.12) is the only in-language extension point, and it
requires the extension to be a **compiled Go plugin**, built with a matching
Go toolchain and OS/arch, and only on Linux/macOS. There is no
pure-Ghost plugin/extension mechanism (e.g., a convention for a `.gs`
file to register itself as reusable library-like code beyond ordinary
`import`).

---

## 11. Existing Functionality — Checklist

A condensed, tickable summary of what §8–§10 specify. Everything here is
already implemented and, per the codebase's own test suite, tested — the
measure of how much of the target in this document the implementation
already meets — except where flagged.

**Language core**
- [x] Dynamic typing, 12 runtime value types plus 3 library-wrapper types
- [x] Integer/float number duality with automatic promotion
- [x] Strings, template literals with nested interpolation
- [x] Lists and maps, with broadcasting arithmetic on lists; `Map` guarantees insertion order (§13.5, §14 decision 2)
- [x] `if`/`else if`/`else`, `while`, C-style `for`, `for ... in`
- [x] `switch`/`case`/`default` as a match-expression (no fallthrough)
- [x] `break`/`continue`/`return`
- [x] First-class functions, closures, default parameters, rest parameters, spread (call sites and list literals), missing-argument checking (§14 decision 1)
- [x] Destructuring assignment (`[a, b] = list`, `{x, y} = map`, `{x: a} = map`)
- [x] Classes, single inheritance, traits/mixins, `this`/`super`
- [x] Per-instance field initialization (no shared-mutable-default bug)
- [x] Module system (`import`, named imports, `import *`, circular-import detection, `scheme:`-prefixed imports)
- [x] Import-only standard library — `console`/`type` global, everything else imported by name (§8.9, §9.1)
- [x] Embedder-claimable import schemes (`RegisterModuleForScheme`/`RegisterFunctionForScheme`, §10.3) — a Go host gets its own `host:name` namespace alongside `ghost:`
- [x] Native classes (`RegisterClassForScheme`, `object.NativeClass`/`Constructible`, §8.9, §10.3) — a Go host can register a class whose instances are built and driven entirely by Go, `new`-able exactly like a Ghost-defined class
- [x] Bounded recursion with a clean error instead of a stack overflow
- [x] Structured, uniform error model with call traces and typo suggestions
- [x] Constant folding and identifier-classification optimization pass
- [x] ANSI-aware, environment-respecting colored diagnostics
- [x] Embeddable in Go, with panic recovery at the boundary

**Standard library**
- [x] `console` — 8 methods
- [x] `math` — ~90 methods/properties across scalar ops, predicates, stats, linear algebra
- [x] `date` — ~35 functions, `date-fns`-modeled, UTC-only
- [x] `random` — seeded PRNG, shared with `math`
- [x] `os` — args/exit/sleep/name
- [x] `file` — read/write/append/exists/isDirectory/size/delete/mkdir/list/copy/move, script-relative paths
- [x] `path` — join/basename/dirname/extname/isAbsolute
- [x] `json` — encode/decode
- [x] `http` — minimal handler/listen server
- [x] `ghost` — abort/execute/extend/identifiers/version (reflection + plugin loading)
- [x] Built-in methods on `string`, `number`, `list`, `map`, `boolean`, `null`, `date`

**Tooling**
- [x] REPL with history and line editing
- [x] File execution with exit-status propagation
- [x] `-h`/`-v`/`-t`/`-i` CLI flags

---

## 12. Required for 1.0: Functionality to Add

The following close the gap between the target this specification sets and
what the implementation does today. Organized by area, in roughly descending
order of how likely each is to be hit by an ordinary user.

**String methods — done.** `contains`, `indexOf`/`lastIndexOf`, `repeat(n)`,
`padStart`/`padEnd`, `charAt`, `slice`, `reverse`, `isEmpty` are implemented
(`object/string.go`, tested in `evaluator/string_methods_test.go`); see §9.1
for the finished reference table and the naming calls made along the way
(`includes`/`at`/`substring` were deliberately not added as second spellings
of `contains`/`charAt`/`slice`).

**List methods — done.** `indexOf`, `find`/`findIndex`, `flatten`, `some`,
`every`, `insertAt`/`removeAt`, `unshift` (front-insert, to pair with
`push`/`pop`/`shift`), `fill`, `chunk`, `flatMap`, `isEmpty` are implemented
(`object/list.go`, tested in `evaluator/list_methods_test.go`); see §9.2 for
the finished reference table and the naming calls made along the way
(`any`/`all` were not added as second spellings of `some`/`every`, and the
generic multi-purpose `splice` was not added at all — `insertAt`/`removeAt`
cover the same ground as two single-purpose methods, matching the existing
`push`/`pop`/`shift` style better than one do-everything call would).

**Map methods — done.** `remove(key)` and `entries()` are implemented
(`object/map.go`, tested in `evaluator/map_methods_test.go`); see §9.2 for
the finished reference table (`remove` was picked over `delete` to match
`list.removeAt()` rather than add a second spelling for the same
operation).

**Number methods — done.** `ceil()`, `abs()`, `pow(exponent)`,
`clamp(low, high)`, and the `isX` predicates (`isNaN`, `isFinite`,
`isInfinite`, `isInteger`, `isEven`, `isOdd`, `isNegative`, `isPositive`,
`isZero`) are implemented as instance methods (`object/number.go`, tested
in `evaluator/number_methods_test.go`); see §9.2 for the finished
reference table. Full parity with `math` was picked over the minimal
`ceil()`-only fix, since `math`'s own predicates and bounds/interpolation
methods (`sqrt()`, `lerp()`, etc.) have no natural instance-method reading
the way `abs`/`pow`/`clamp`/`isX` do — those stay `math`-only.

**Date module — done.** `subHours`/`subMinutes`/`subSeconds` and
`startOfWeek`/`endOfWeek`/`startOfYear`/`endOfYear` are implemented
(`library/modules/date.go`, tested in `library/modules/date_test.go`); see
§9.5 for the finished reference (`startOfWeek`/`endOfWeek` treat Sunday as
the first day, matching `weekday()`'s existing `0 = Sunday` convention and
`date-fns`'s own default, rather than an ISO 8601 Monday-first week).

**Language-level gaps — done.**
- **Destructuring assignment** (`[a, b] = list`, `{x, y} = map`, and `{x: a}
  = map` to bind under a different name) is implemented (`ast/destructure.go`,
  `parser/destructure.go`, `evaluator/assign.go`; tested in
  `parser/parser_test.go`'s `TestListPatternAssignment`/`TestMapPatternAssignment`
  and `evaluator/evaluator_test.go`'s `TestListPatternAssignment`/
  `TestMapPatternAssignment`); see §8.3 for the finished reference and its
  restrictions (statement-level only, plain names only, no nesting).
- **Rest parameters and spread syntax**, both directions, are implemented:
  a function's last parameter may be `...name` (collects the rest as a
  list), and `...expr` expands a list's elements in place at a call site or
  in a list literal (`ast/spread.go`, `parser/spread.go`,
  `evaluator/expressions.go`, `evaluator/function.go`; tested in
  `evaluator/evaluator_test.go`'s `TestRestParameters`/`TestSpreadExpressions`
  and `parser/parser_test.go`'s `TestFunctionRestParameter`/
  `TestSpreadExpression`); see §8.7 for the finished reference.
- **Arity checking for user-defined functions/methods** — see §14 decision 1
  (done).

**Arithmetic gaps found building a library on Ghost.** Two small, concrete
additions, from the same Chisel/Studio work that produced §13.13–§13.20:

- **Floor division.** Ghost has true division only: `7 / 2` is `3.5`, and
  `6 / 3` is `2` (a `number`, equal to `2` — there is no spurious float
  promotion, and the int/float rules in §8.4 are working as documented).
  What is missing is a way to *ask* for the integer quotient. Pixel layout
  is the motivating case: every grid, cell, scroll offset and zoom step
  wants `floor(a / b)`, and writing that composition on every line is a real
  tax — Chisel's `Rect` exists partly to pay it once and hand back whole
  numbers. Add `math.floorDiv(a, b)`, matching the existing `math`
  vocabulary. Note the obvious operator spelling is unavailable: `//`
  is a comment.
- **`%=`.** `+=`, `-=`, `*=` and `/=` all work; `%=` alone is a syntax error
  ("`=` cannot start an expression"), even though `%` is a first-class
  arithmetic operator in §8.4 and `compound.go` handles the other four. This
  is an omission in the compound-operator set rather than a design stance,
  and closing it is a table entry, not a feature.

---

## 13. Required for 1.0: Defects to Fix

Ranked roughly by how much damage each can do, not by how easy it is to fix.
Every finding below was confirmed by reading the relevant code paths
directly (file/line references given); none are speculative.

### 13.1 The shared random-number generator is not safe for concurrent use — done

`library/modules/random.go`'s `seed` and `randomizer` (`*rand.Rand`) are now
behind a package-level `randomState sync.Mutex`, mirroring `moduleState` in
`evaluator/import.go`. Every read or write — `randomRandom`'s and
`mathRandomInt`'s draws (via the new `randomFloat64`/`randomInt63n`
accessors), `randomCurrentSeed`'s read, and `SeedRandom`'s write (reached by
both `randomSeed` and `mathRandomSeed`) — takes the lock for the whole
operation, since Go's `*rand.Rand` is not safe for concurrent use on its own
and a lock held only around reading the pointer would not protect the
`Float64()`/`Int63n()` call's own internal state mutation. Tested in
`library/modules/random_test.go`'s `TestRandomConcurrentAccessIsRaceFree`,
which spins up 100 pairs of goroutines calling `random.random()` and
`math.randomInt()` concurrently — confirmed to fail under `go test -race`
against the unguarded code (a real data race inside `math/rand`, not a
false positive) and to pass clean, repeatedly, against the fix.

### 13.2 `==`/`!=` cannot compare two maps, functions, or several other same-typed pairs — not even by identity — done

`evaluator/infix.go`'s `evaluateEquality` used to special-case only `NULL`,
`INSTANCE` (identity), `LIST` (deep), and `DURATION` (deep); everything
else — same-typed `MAP`, `FUNCTION`, `CLASS`, `TRAIT`, `SUPER`, `SCOPE`, the
`LIBRARY_*` wrapper types — fell through to `evaluateInfix`'s final
`operatorError`, producing `` cannot use `==` between two maps `` even for
`m == m` (the identical object compared to itself). This directly
contradicted `object/equality.go`'s own doc comment on `ValuesEqual` —
*"it is what `==` means between two Ghost values... a value counted as
equal inside a list is equal everywhere else in the language too"* — since
`object.ValuesEqual` (used by `list.contains()`/`list.unique()`) had no
`Map` case either and fell back to identity for it, while `evaluateInfix`
never called `object.ValuesEqual` at all — it had its own, narrower,
hand-rolled `evaluateEquality`, so the two had silently diverged in two
directions at once.

**Fix:** `evaluateEquality` now routes every comparison except a same-typed
pair of `Boolean`/`Number`/`String`/`Date` (each of which keeps reading its
own dedicated infix evaluator, for a reason specific to that type — see
§8.5) through `object.ValuesEqual` directly, which is now genuinely the one
place equality is decided. `ValuesEqual` itself gained the missing `Map`
case (`MapsEqual`, content equality mirroring `ListsEqual`) and a `Date`
case (instant equality, matching `evaluateDateInfix`'s existing `==` — this
also fixes the same divergence for `list.contains()`/`unique()` on a list
of dates, which previously used identity despite `==` between two dates
using the instant); `Duration`'s existing field-by-field comparison moved
from the evaluator into `object.DurationsEqual` alongside it. Everything
still unhandled by a specific case (`Function`, `Class`, `Trait`, `Super`,
`Scope`, the library wrapper types, and any type an embedding host defines
of its own) now falls back to identity — a real answer instead of a type
error — the same rule `Instance` already had. Tested in
`object/equality_test.go` (`ValuesEqual` directly, including the identity
fallback via a minimal stub type) and `evaluator/evaluator_test.go`'s
`TestEqualityComparisons`/`TestEqualityTypeMismatch` (confirming a
same-type comparison now answers rather than errors, and a genuinely
different-typed, non-null pair still errors exactly as before).

### 13.3 `switch` silently swallows an error in its subject or case expressions — done

`evaluator/switch.go`'s `evaluateSwitch` used to call `Evaluate(node.Value,
scope)` and, for each case, `Evaluate(val, scope)`, without ever checking
`isError()` on either result: if the switch's subject (or a case value)
failed to evaluate, the resulting `*object.Error` was compared via `.Type()`/
`.String()` against the other branches instead of being propagated — in the
best case falling through to `default` (or returning `nil` if there was
none), silently discarding a real error instead of surfacing it the way
every other construct in the language does.

Separately, case-value comparison used `obj.Type() == out.Type() &&
obj.String() == out.String()` rather than `object.ValuesEqual` — a
comparison by *string representation*, which was not just weaker than the
rest of the language's equality rules (§8.5) but actively wrong for two
whole types: `Function.String()` answers the literal `"function"` for
*every* function, and `Class.String()` answers `"class"` for every class —
neither type has a more specific string form — so `switch (f) {
case g { ... } }` matched on **any** two functions regardless of which ones
they actually were, and the same for two classes.

**Fix:** `evaluateSwitch` now checks `isError()` immediately after
evaluating the subject and each case value and propagates it, and compares
every case through `object.ValuesEqual` (§13.2) instead of `String()` —
content equality for a list/map/duration/date, identity for everything
else, matching `==` exactly. Tested in `evaluator/switch_test.go`
(`TestSwitchPropagatesErrors`, and `TestSwitchComparesByValueNotString`,
which reproduces the function/class false-match directly).

### 13.4 `-i` is documented in two places and implemented in none — done

Both `README.md` ("Interactive mode... pass the `-i` flag") and
`cmd/help.go`'s own `-h` output ("`-i` enter interactive mode after
executing file") described a flag that `cmd/ghost.go` never registered
(`flag.BoolVar` was called for `flagHelp`, `flagVersion`, `flagTime` only),
and there was no code path anywhere that ran a file and then started a REPL
with its environment intact. Confirmed before the fix: an unregistered `-i`
does not silently no-op — Go's `flag` package rejects it and `flag.Parse()`
exits the process with its own "flag provided but not defined" error, which
is a worse failure mode than the ignored-flag one the original writeup
guessed at.

**Fix:** `-i` is now registered in `cmd/ghost.go` and, when a file argument
is also given, runs the file exactly as before and then hands the same
`*ghost.Ghost` instance — same `Scope`, so every binding the script made is
still there, even one made before a mid-script failure the script's own
error report already printed — to a new `repl.StartWithInstance` (the
existing `repl.Start` now just builds a fresh instance and calls it),
rather than the REPL's usual brand-new one. The failed-script exit status
(`os.Exit(1)`) is skipped once `-i` hands off to the REPL; the interactive
session decides its own exit the way a bare `ghost` invocation already
does. `-i` with no file argument is a no-op — falls through to the same
"start a plain REPL" branch as no arguments at all — since there is no file
to execute first. Tested in `ghost/ghost_test.go`'s
`TestExecuteReusesScopeAcrossCalls` (the scope-continuity mechanism itself;
the REPL and CLI layers around it need a terminal to drive and are outside
this codebase's existing test boundary — `repl/` and `cmd/` have no test
files today, and this doesn't change that) and confirmed directly against
the built binary (piped stdin) for every row `-i` interacts with: no file,
`-t` composed with `-i`, and a script that fails partway through.

### 13.5 Map/list `for ... in` and `keys()`/`values()` iterate in random order — done

`evaluator/for_in.go`'s map branch and `object.Map`'s `keys()`/`values()`/
`String()` used to range directly over the underlying Go `map[MapKey]MapPair`,
whose iteration order Go deliberately randomizes — not just between two
runs of the identical script, but between two calls in the very same run,
on the very same `Map`. This was inconsistent with the rest of the
language's posture toward predictability (`list.sort()` is explicit and
stable; the "did you mean" suggestions are tie-broken deterministically)
and was a common source of "flaky" script output and non-reproducible test
assertions for users.

**Fix:** `object.Map` (`object/map.go`) now carries its Go map (`Pairs`,
still exported and untouched for O(1) lookup by `get`/`has`/`[]`/`.`) beside
an unexported `order []MapKey` slice that only its own `SetPair`/
`RemovePair` methods touch — a new key is appended, an existing key keeps
its original position and only its value changes, and a removed key leaves
no gap. `OrderedPairs()` is the one place that order is read back, so
`keys()`, `values()`, `entries()`, `String()`, and `for ... in` all agree
with each other and with `SetPair` on what it is. Every mutation in the
language — a map literal, `set()`, index/property assignment, `merge()`
(this map's pairs first in their own order, then the other's remaining
pairs in its order, so a shared key keeps this map's position for it and
takes the other's value — the same result a plain object spread
(`{...left, ...right}`) gives in JS) — goes through `SetPair`/`RemovePair`,
never `Pairs` directly, so `order` can't drift out of sync with it.

One piece this doesn't reach: `json.encode()`'s own `encoding/json.Marshal`
call always alphabetizes a Go map's keys regardless of the order fed into it,
so there was never an order for `json.encode()` to preserve or lose in the
first place; `json.decode()` similarly can't recover a source JSON text's
key order from `encoding/json.Unmarshal`'s own generic-map decoding. Both
still settle into *a* fixed, repeatable order once decoded/loaded (per this
fix), just not one that means "the order written in the JSON text" the way
a Ghost map literal's does mean "the order written in the source." A
`Map` built by an embedding host from a plain Go map (`object.NewMap`,
`AnyValueToObject`'s map case) has the same limitation, for the same
reason: a Go map carries no order of its own to preserve.

The AST needed the same fix first: `ast.Map.Pairs` was itself a Go map
(`map[ExpressionNode]ExpressionNode`) keyed by the key expression, which
lost a map literal's source order before evaluation ever got the chance
to preserve it. It is now `[]MapEntry` (`ast/map.go`), ordered the way it
was written; `parser/map.go` appends to it instead of indexing into it,
and `evaluator/map.go`, `parser/destructure.go`, and `optimizer/optimizer.go`
read it as a slice. A repeated key in a literal (`{x: 1, x: 2}`) now
deterministically keeps the position of its first appearance with the
last value written to it, the same rule `set()` follows, rather than it
being arbitrary which of the two occurrences won.

Tested in `object/map_test.go` (`SetPair`/`RemovePair`/`OrderedPairs`
directly), `evaluator/map_methods_test.go`'s `TestMapPreservesInsertionOrder`
/`TestMapForInPreservesInsertionOrder`, and `parser/parser_test.go`'s
`TestMapLiteralPairsPreserveSourceOrder`.

### 13.6 Bounds-checking is inconsistent across list/map/string operations — done

`list[i]`/`map[k]`/`string[i]` indexing all return `null` for an
out-of-range index or missing key (`evaluator/index.go`), as do
`list.pop()`/`shift()`/`first()`/`last()`/`tail()` on an empty list — but
`list.slice()` raises an `Index` error for an out-of-range `start`/`end`
(`object/list.go`). Neither behavior is wrong on its own, but the split was
not documented as an intentional rule anywhere, which made it easy for a
new list/string method to pick the "wrong" one of the two conventions by
accident.

**Fix:** the rule is now stated explicitly — "a read that names a *position*
is lenient; an operation that names a *range* validates it" — in §9.2's
intro, immediately above the per-type method tables it governs, rather than
in §7: §7 is scoped to *naming*, by its own opening sentence, and this is a
behavioral convention, so folding it in there would be off-topic for a
reader looking either up. §9.2 is where every method affected by the rule is
already documented, and several of its table rows (`charAt()`, `slice()`,
`insertAt()`/`removeAt()`) already pointed back at "§13.6" for their
reasoning before this fix existed anywhere to point to — they now point at
the rule itself.

The audit turned up one real bug, not just an undocumented split:
`evaluator/index.go`'s `evaluateStringIndex` (the `string[i]` operator)
checked `idx` against `len(str.Value) - 1` — a **byte** count — while
`object/string.go`'s `charAt()`/`slice()`/`length()` all correctly use a
rune count (`utf8.RuneCountInString`/`[]rune`), per §"Working through
SPEC.md §11–14" in `CLAUDE.md`. For a receiver with any multi-byte rune,
byte count exceeds rune count, so an index in between passed the (too
generous) bounds check and then indexed past the end of the `[]rune`
conversion that came after it — a genuine Go panic, not just a wrong answer,
recovered by `ghost.Execute` into an unhelpful `internal error` instead of
the `null` a position read out of range is supposed to give. Confirmed
before the fix: `"héllo"[5]` (5 runes, 6 bytes) crashed; `"héllo"[4]`
(the last valid rune) did not, since it stayed under both counts. Every
other list/map/string bounds check audited against the rule already matched
it correctly: list/map/string indexing and `list`'s empty-collection reads
were already lenient; `list.slice()`/`string.slice()`/`list.fill()` were
already validating; `list.insertAt()`'s clamp-to-nearest-end is a
deliberately different leniency, justified in its own §9.2 row, for a
write that has no "absent" position to fall back to the way a read does.

Tested in `evaluator/index_test.go`'s `TestIndexingIsLenientOnPosition`
(every position-reading indexing form, list/map/string, answers `null`
rather than erroring — including the exact multi-byte case that used to
panic) and `TestStringIndexUsesRunePositions` (confirms `string[i]` and
`charAt()`/`length()` agree on rune positions for multi-byte input).

### 13.7 `string.find`/`findAll`/`matches` invert the usual receiver/argument relationship — done

`(pattern string).find(subject string)` used to compile the **receiver** as
the regular expression and match it against the **argument**. Every widely
used language Ghost is modeling itself on does the opposite —
JS's `subject.match(pattern)`, PHP's `preg_match(pattern, subject)` (subject
still second, but the *call* reads pattern-first, not method-chained on the
subject), Python's `re.search(pattern, subject)`. A Ghost user writing
`text.find(pattern)` out of that muscle memory got a working call that did
something different from what they meant, with no error to catch the
mistake (both operands are strings). This was called out and clearly
intentional in the code's own comment ("the string a pattern method is
called on is the pattern itself"), but §14 decision 3 called for it to flip
to `subject.find(pattern)` before 1.0 locked the API, while it was still a
breaking change nobody had shipped against yet.

Separately, `findAll` was misleadingly named — it called
`FindStringSubmatch` (the *first* match's capture groups), not
`FindAllString`/`FindAllStringSubmatch`. A user reaching for "all the
matches in this text" via `findAll` got only the first match's submatches.
This was a genuine bug, fixed at the same time as the receiver/argument flip.

**Fix:** `object/string.go`'s `find`/`findAll`/`matches` now all read their
pattern from the *argument*, compiling it fresh on every call
(`compilePattern`, a free function — it no longer has anything receiver-
specific about it, unlike the `(str *String) pattern(...)` method it
replaces), and match it against the receiver. `find` now calls
`Regexp.FindString` directly (equivalent to the old code's `found[0]`, just
without a detour through submatches it never used); `findAll` now calls
`Regexp.FindAllString(receiver, -1)`, answering every match's full text, in
order, rather than one match's capture groups. `matches` is unchanged in
shape beyond the same argument/receiver swap. The bad-pattern error message's
help text was reworded to match ("the argument to a pattern method is the
pattern; call it on the string you want to search").

This is a breaking change to the two call sites shipped in-repo:
`examples/ada.gs` and `examples/ada/ada.gs` (an ELIZA-style chatbot demo)
both called `knowledge.pattern.matches(text)`/`.findAll(text)` in the old
shape and now call `text.matches(knowledge.pattern)`/`.findAll(knowledge.pattern)`.
Flipping the call shape alone was not enough to keep that demo's behavior
intact, though: its response templates (`examples/ada/modules/therapist.gs`)
substitute `{1}`-style placeholders with a match's *capture group*, which
depended on `findAll`'s bug (the accidental capture-group access) as its only
way to reach one — Ghost's string API has never had a documented way to read
a capture group, and adding one is a separate feature, not part of this
defect's scope. Both example files now carry a comment noting the gap rather
than silently shipping a demo whose reflection substitution quietly stopped
working. Neither example is exercised by the test suite (`examples/` has no
Go test referencing it), so this doesn't show up as a regression there.

Tested in `evaluator/string_methods_test.go`'s `TestStringPatternMethods`
(the flipped shape for `find`/`findAll`, `findAll` returning every match
rather than one match's submatches) and the `matches` cases added to
`TestStringMethodBooleans`; `TestBadPatternsAreReported`
(`evaluator/errors_test.go`) was updated to raise its bad pattern through the
argument instead of the receiver, matching the new shape.

### 13.8 `ghost.extend()` is unavailable on Windows

Go's `plugin` package (used by `ghost.extend`, §9.12) only supports
Linux and macOS; there is no Windows implementation at all (it's a build
failure/`ErrNotSupported` situation, not a degraded one). Ghost itself
builds and ships Windows binaries (`Makefile`, `.goreleaser.yml`). The only
in-language way to extend Ghost with native code is therefore unavailable
on a third of Ghost's own supported platforms, with no fallback and no
documented caveat. At minimum, this needs a documented caveat before 1.0;
a real fallback is a larger question outside this specification's scope.

### 13.9 Dead token: `token.PRINT` — done

`token/token.go` defined `PRINT` as a token type, but the scanner's
`keywords` map (`scanner/scanner.go`) never mapped the string `"print"` to
it — so this token could never actually be produced by the scanner, and
nothing in the parser or evaluator referenced it either. Ghost has no
`print` at all in 1.0 (§9.1) — `console.log(...)` replaced it — so there was
no live identifier this token could even be mistaken for. `PRINT` was pure
dead code, presumably left over from an earlier design. **Fix:** deleted,
along with its `typeNames` entry (`token/token.go`).

### 13.10 `-t` is implemented but undocumented — done

The reverse of §13.4: `cmd/ghost.go` registers and honors `-t` ("display
how long the program ran for"), but `cmd/help.go`'s `helpCommand()` never
mentioned it in its printed usage/flags list. **Fix:** added to the printed
help, matching the wording of its own `flag.BoolVar` description
(`cmd/help.go`).

### 13.11 Minor: duplicated `isTruthy` logic — done

`evaluator/evaluator.go` and `object/boolean.go` each defined their own
private `isTruthy(Object) bool` with identical logic (`object/boolean.go`'s
copy is exported as `IsTrue`/`IsFalse`). Framed in §13's original writeup as
"not a bug — both copies agree today" — but the audit this triggered found a
**third**, independent copy, and it had already drifted: `evaluator/prefix.go`'s
`BANG` (`!`) case hand-rolled its own truthiness switch rather than calling
either existing one, and its `default` branch answered `false` for every
type but `Boolean`/`Null` — silently skipping the case `String` has in the
real rule (an empty string is falsy, §8.5). `!""` answered `false` instead
of `true`. A **fourth** copy, in the constant-folding optimizer
(`optimizer/fold.go`'s `foldPrefix`), had the identical bug for the same
reason — it explicitly commented that it was "matching the evaluator," and
did, including the bug, folding every string literal's `!` to `false`
regardless of content. Confirmed before the fix: `!""` (and `!("a".slice(0,
0))`, to rule out the optimizer folding a literal masking a runtime-path
bug) both answered `false` against the built binary.

**Fix:** `evaluator/evaluator.go`'s private `isTruthy` is deleted; its four
call sites (`evaluator/if.go`, `ternary.go`, `while.go`, `for.go`) now call
`object.IsTrue` directly, so `object/boolean.go`'s `isTruthy` (reached only
through the exported `IsTrue`/`IsFalse`) is genuinely the one place this
rule is decided, the same principle §"Error handling" already applies to
`object.ValuesEqual` (§13.2). `evaluator/prefix.go`'s `BANG` case is
replaced with `toBooleanValue(object.IsFalse(right))` — no more hand-rolled
switch. `optimizer/fold.go`'s `foldPrefix` keeps its own switch (folding
happens over `ast` nodes, before any `object.Object` exists to call
`IsFalse` on), but its `String` case now checks `right.Value == ""` instead
of unconditionally folding to `false`; its `Number` case is unchanged and
correctly unconditional, since a number is always truthy regardless of
value (§8.5) — `!0` really is `false`.

Tested in `evaluator/evaluator_test.go`'s `TestBangOperator` (new `!""` /
`!0` cases, and reworded comments explaining *why* each answer is what it
is rather than restating the old, incomplete "non-boolean, non-null
operands remain falsy" framing that this bug slipped through) and
`optimizer/optimizer_test.go`'s `TestFoldsComparisonsAndBooleans` (the same
`!""` case, at the constant-folding layer specifically, so a regression in
either the runtime or the compile-time path is caught independently).

### 13.12 Two statements on separate lines can silently merge into one — done (partially, by design)

§8.1's own "statement separation" bullet used to claim "there is no
significant-newline rule and no automatic-semicolon-insertion logic to
reason about" — false in two directions, both confirmed by running the CLI
directly (not just the evaluator test helper, which never checks
`parser.Errors()` and so does not catch either):

1. **A statement can swallow the next line's opener.** `parseExpression`'s
   Pratt loop (`parser/expression.go`) only stops at a `;` or a token with no
   infix meaning; a newline is not significant, so a completed statement
   followed on the next line by `[`, `(`, `.`, `++`/`--`, or a binary
   operator continues parsing as if the two were one expression. `x = 1\n[10,
   20, 30]` parses as `x = 1[10, 20, 30]` (an index expression), not two
   statements — confirmed to fail with `type error: cannot index number` for
   a single-element list and a comma-related syntax error for more than one.
   Ending the first line with an explicit `;` avoids it (`assign()` stops the
   value expression at the `;`), which is presumably why this has not been
   noticed: idiomatic Ghost code chains same-line statements with `;` and
   rarely starts a line with `[`/`(` right after an unrelated one.
2. **A `;` after anything but an assignment breaks the next statement.**
   `assign()` explicitly consumes a trailing `;` (`parser/assign.go`), but
   `expressionStatement()` (`parser/statement.go`) does not — so a bare
   expression statement (a call, `console.log(1)`, included) followed by
   `;` leaves the block/program loop re-entering `statement()` with the `;`
   itself as the current token, and there is no prefix parser for
   `SEMICOLON`. `console.log(1);\nconsole.log(2);` fails to parse at all
   ("`;` cannot start an expression"), even though every other call-heavy
   example in this document is written exactly that way.

The original writeup offered two alternative fixes: make statement
termination actually newline-significant, or give every statement-producing
path the same trailing-`;`-consumption `assign()` already has, symmetrically.
These fix different bugs — newline-significance would resolve both; the
symmetric-consumption path resolves only #2 outright, leaving #1's
already-rare, `;`-avoidable case as an accepted, documented trade-off rather
than a fixed bug. The narrower fix was chosen deliberately: full
newline-significance is a materially larger grammar change whose own
correct behavior is not obvious in every case — Ghost's class/method syntax
already invites JS-style fluent chaining (`foo()\n  .bar()`), which
newline-significant termination would need to special-case correctly rather
than break, and nothing in this codebase today establishes whether that
idiom is meant to be supported. Symmetric `;` consumption carries no such
risk: it only ever makes a previously-rejected program (one correctly
terminated with `;`) parse, never changes what an already-accepted program
means.

**Fix (bug #2):** `parser/statement.go`'s `expressionStatement()` — the path
every bare expression statement, and every `if`/`while`/`for`/`function`/
`class`/`trait`/`switch`/`import`/`use`/`break`/`continue` statement (each a
prefix parser reached through `parseExpression`, per `parser.go`'s
`registerPrefix` table), funnels through — now consumes an optional
trailing `;` the same way `assign()`/`returnStatement()` already did. This
alone did not cover a bare list/map-literal statement (`[1, 2, 3]`,
`{"a": 1}`): `assign()` special-cases anything starting with `[`/`{` into
`destructuringAssign()` before `expressionStatement()` is ever reached, to
read it first as a possible destructuring pattern (§12); each of
`destructuringAssign()`'s "this wasn't actually a pattern" fallback
branches built its own `&ast.Expression{...}` node directly, bypassing
`expressionStatement()` (and its now-fixed trailing-`;` consumption)
entirely. Both functions' "wrap this expression as a statement" logic is
now the single shared `expressionStatementFrom` (`parser/statement.go`),
called from `expressionStatement()` and from all four of
`destructuringAssign()`'s fallback returns, so there is exactly one place
this is decided rather than two that happened to agree until one of them
didn't.

**Bug #1 stays as documented, existing behavior** — §8.1's "statement
separation" bullet is reworded to describe it accurately instead of denying
it exists.

Tested in `parser/parser_test.go`'s `TestSemicolonTerminatesEveryStatementKind`,
which runs each of the 14 statement kinds above (plus both destructuring
forms, to confirm the fallback-path fix specifically) through
`parser.Errors()` directly, exactly the multi-statement/multi-line-through-
`Errors()` check this callout asked for — confirmed to fail on all but the
already-working assignment/return cases against the pre-fix code, and pass
clean against the fix.

### Findings 13.13–13.20: from building a library on Ghost

The callouts below came out of writing Chisel (a retained-mode widget
toolkit) and Studio (a sprite and tilemap editor shell built on it) against
Ghost 1.0 — the first substantial body of Ghost code that is a *library*
rather than a script, which is why they cluster around scoping, module
boundaries, and class shape rather than around the standard-library surface
§12 covers. Every behavior below was executed against a build of the
interpreter rather than read off this document; where the two disagree, the
disagreement is itself the callout, and §13.15 and §13.17 are both cases
where this specification already describes the behavior we want and the
interpreter does not implement it.

They are appended rather than interleaved into the damage ranking above:
§13.1–§13.12 are closed, and renumbering them would invalidate every
cross-reference in this document and in the commit history citing them.
Within this block the ranking convention still holds. §13.13, §13.14 and
§13.15 were one design question wearing three faces; §14 decision 9 settled
it and all three were fixed together, since fixing any one in isolation would
have changed what the other two meant.

### 13.13 A function cannot assign to a variable outside itself — done

```ghost
scale = 1
function set(n) { scale = n }
set(4)
console.log(scale)   // was 1; now 4
```

`evaluator/assign.go`'s `evaluateIdentifierAssignment` wrote through
`object.Environment.Set`, which by its own doc comment never walked the outer
chain. Reads were not symmetric with writes: `Get` and `Has` recursed through
`outer`, `Set` did not. A function could therefore *see* an outer name, and
mutate an object held under it, but never rebind it — the assignment silently
created a frame-local binding that died with the call, with nothing to warn at
the write and nothing to raise at the stale read afterwards.

**Fix:** `object.Environment` gains `Assign`, which walks the enclosing chain
exactly as `Get` does and rebinds the name wherever it is already bound,
reporting whether it found one. It is built on a new `rebind` helper that
updates an existing binding without ever creating one, which `Set` now uses
too, so the "update in place" logic exists once rather than twice. The
evaluator's `bind` (`evaluator/assign.go`) is the single point every
name-assignment goes through — `Assign` first, falling back to `Set` to
declare locally when the name is bound nowhere. Plain assignment, both
destructuring forms (`evaluateListPatternAssignment`,
`evaluateMapPatternAssignment`), compound assignment (`evaluator/compound.go`)
and `++`/`--` (`evaluator/postfix.go`) all route through it, so every spelling
of assignment reaches an outer variable the same way. §8.3 now documents the
rule.

The cost §14 decision 9 weighed is real and now live: a local whose name
matches a sibling method rebinds that method, since a method body's scope is
the class environment. Block scoping (§13.15) is what contains it — a name
first assigned inside a block stays there — and
`TestMethodScopingIsUnchanged` pins the three cases that must not blur (a
method's local does not escape into the instance, a method does rebind a
module-level variable, and a field assignment still reaches the instance).

Tested in `evaluator/scoping_test.go`'s
`TestAssignmentReachesAnEnclosingScope`, which covers the direct case, a
nested function reaching two levels out, a parameter correctly shadowing
rather than rebinding, a genuinely new name staying local, and each of the
compound/postfix/destructuring spellings.

### 13.14 Closures created in a loop cannot capture the loop variable — done

```ghost
handlers = []
for (name in ["a", "b"]) {
    handlers.push(function () { return name })
}
handlers[0]()   // was: name error: `name` is not defined; now "a"
```

`evaluator/for_in.go` bound the loop's control variables with
`scope.Environment.Set` into the *enclosing* environment and `Delete`d them in
a deferred restore once the loop ended (`evaluator/for.go` did the same for a
C-style loop's identifier). A closure created in the body captured that
enclosing environment by reference rather than the value the variable held on
its iteration, so every closure shared one binding — and once the restore had
run, no binding at all. The error surfaced at call time, arbitrarily far from
the loop, and named the innocent variable.

**Fix:** both loops now bind their control variables in an environment created
for that iteration alone, and the save-and-restore is gone entirely — the
variables were never written to the enclosing scope, so there is nothing to
put back. `evaluateForInBody` (`evaluator/for_in.go`) runs one iteration in
`enclose(scope)` with the key and value set there. `evaluateFor`
(`evaluator/for.go`) carries the control variable between iterations in a Go
local and re-declares it in each iteration's scope, which also let the loop
drop to one environment level rather than two.

One ordering detail is load-bearing: the increment runs at the *top* of each
iteration after the first, against that iteration's own scope, rather than at
the foot of the previous one. Incrementing in the iteration a closure just
captured would move the value out from under it, which is the same bug in a
new place.

Tested in `evaluator/scoping_test.go`'s `TestClosuresCaptureTheirIteration`
(both loop forms, a `while` body's local, a closure made inside a nested block
inside a loop, and a closure that keeps mutating its own captured binding
across calls) and `TestLoopVariablesDoNotDisturbTheEnclosingScope` (the loop
variable does not overwrite a same-named variable outside it, and accumulators
declared outside the loop still accumulate).

### 13.15 Blocks do not introduce a scope — done

```ghost
x = 1
if (true) { x = 2; y = 99 }
console.log(x)   // 2 — assignment still reaches the outer x
console.log(y)   // was 99; now: name error: `y` is not defined
```

§8.3 stated that blocks each introduce a new enclosed `Environment`, and only
function bodies did. `evaluator/block.go`'s `evaluateBlock` threaded the
caller's scope straight through, so an assignment inside an `if`, `switch`,
`while` or `for` body wrote to the enclosing function scope.

**Fix:** the scope is introduced by the statements that *own* a block —
`evaluateBranch` (`evaluator/if.go`), `evaluateWhile`, `evaluateCase`
(`evaluator/switch.go`), and both loops — rather than by `evaluateBlock`
itself. That distinction matters: two of `evaluateBlock`'s callers must not
get a scope, since a class or trait body is evaluated directly in the
environment that collects its members, and a function or method body already
runs in the frame `createFunctionEnvironment` built for it.

**This is a breaking change**, and the only one in this block of callouts. A
name first assigned inside a branch no longer outlives it, so
`if (c) { result = 1 } else { result = 2 }` followed by a read of `result` is
now a `Name` fault where it used to work. Assigning `result` before the branch
fixes it, and that is the form §8.3 now shows. Every one of the 41 programs in
`examples/` produces byte-identical output before and after, so the pattern is
rarer in practice than it looks.

**Performance.** A scope per block execution is an allocation per loop
iteration, which cost up to 8.5× the allocated bytes and 1.8× the wall time on
`evaluator/benchmark_test.go`'s loop-heavy cases. Two changes bring it back:
`object.Scope.Enclose`/`Release` keep one finished block scope per environment
(`Environment.freeChild`) and hand it to the next block rather than allocating,
and `Environment.Capture` marks an environment — and the whole chain enclosing
it — whenever a closure, class, or trait is created inside it, so a captured
scope is dropped instead of reused and the value that captured it keeps
reading what it closed over. Reuse is safe across goroutines by construction:
every concurrent entry into Ghost code (an `http.handle` callback, an
embedder's `Call`) runs in a function frame of its own, so a block's
environment is a child of that frame rather than of anything shared, and
`go test -race` passes. Allocation is now at parity with the pre-change
interpreter on every benchmark; wall time is 1.03×–1.21×, which is the
intrinsic cost of the extra link every name lookup crosses.

Tested in `evaluator/scoping_test.go`'s `TestBlocksIntroduceAScope` (a name
first assigned in an `if`, `else`, `while`, `for` or `switch` case body does
not outlive it, and neither loop's control variable outlives its loop), with
the reuse machinery's correctness covered by
`TestClosuresCaptureTheirIteration` — every case there would fail if a
captured environment were handed to the next iteration.

### 13.16 Every top-level name in a module is exported, including its own imports

```ghost
// mod/lib.gs
import "ghost:math" as math
helper = "PRIVATE"
class Public { }
```
```ghost
import "mod/lib" as lib
lib.helper   // "PRIVATE"
lib.math     // library_module — lib's own import, re-exported
```

A module's scope *is* its export surface: `evaluateImport` binds the loaded
module's scope and every name in it is reachable through the alias. A file
therefore cannot keep a private helper beside the public class it supports,
and — the sharper half — importing a module pulls that module's own imports
along with it, so `lib.math` resolves even though `lib` never offered it.
That makes a module's public surface a function of its implementation
details: adding an import to a file silently adds a name to its API, and
renaming a private helper is a breaking change to consumers who were never
supposed to see it.

One-class-per-file makes this survivable rather than solving it, and it is
why Chisel's layout has thirty small files rather than six coherent ones —
the file layout is working around the language, which is the wrong reason
for a file layout.

**Severity: high, design.** Suggested fix: an explicit `export` marker, with
the current export-everything behavior retained as the fallback for a file
that declares no exports at all, so existing code keeps working unchanged
and only a file that opts in gets a curated surface.

### 13.17 Calling a sibling method by bare name loses the receiver — §8.8 says it works

```ghost
class Widget {
    constructor(n) { this.n = n }
    describe() { return this.n }
    show() { return describe() }     // name error: `this` can only be used inside a class
}
```

§8.8 states that "a method body's scope is the class environment, so sibling
methods call each other by bare name with no `this.` required." The bare
call resolves — the sibling is found in the class environment — but it is
invoked with the wrong receiver, so it fails the moment its body touches
`this`.

The cause is in `evaluator/call.go`'s `unwrapCall`. A resolved
`*object.Function` gets `functionScope.Self = callee` (the function itself),
and the receiver is only carried forward by the branch that re-binds `Self`
when `callee.Scope.Self` is an `*object.Instance`. For a method, the scope
captured at declaration is the *class body's* scope, whose `Self` is the
`*object.Class` — never an `*object.Instance` — so that branch does not
fire, and the method runs with no instance bound. A method reached the
documented way, through `evaluateMethod` → `callInstanceMethod` →
`invokeMethod`, is given `Self: instance` explicitly, which is why
`this.describe()` works and `describe()` does not.

Two follow-on notes: a sibling method that happens *not* to touch `this`
works fine, so the failure depends on the callee's body rather than on the
call — which is how it survived to 1.0. And the error is reported at the
callee's `this`, not at the bare call that supplied the wrong receiver, so
it points at correct code and hides the actual mistake.

**Severity: mid, documentation drift.** Independent of §14 decision 9 and
fixable on its own: either carry the receiver through the bare-call path, or
retract §8.8's claim and require `this.`.

### 13.18 A field and a method may share one name, with no diagnostic

```ghost
class Thing {
    thing = 7
    thing() { return "method" }
}
t = new Thing()
t.thing     // 7
t.thing()   // "method"
```

Fields are initialized into `instance.Environment` (`initializeField`);
methods live in `class.Environment`. A property *read* goes through
`evaluateInstanceProperty`, which consults `instance.Environment.GetLocal`
first and so answers with the field; a *call* goes through
`evaluateMethod` → `callInstanceMethod`, which resolves against the class
chain and so finds the method. Neither path knows the other exists, so the
two never collide and nothing is ever reported.

One name meaning two different things depending on whether it is followed by
parentheses is not a behavior any reader will predict, and the declaration
that creates it is silent at the point where it could most cheaply be
caught. Note that this is the *opposite* of what a reasonable person
assumes on reading the class — the assumption that the field makes the
method unreachable is itself enough to misdesign around, which is how this
was found.

**Severity: mid.** Suggested fix: reject the duplicate declaration at class
construction with a `Syntax` fault naming both, rather than changing either
lookup path.

### 13.19 Module resolution is a global, first-match-wins search path

`evaluator/import.go`'s `resolveModule` calls `addSearchPath` with the
importing file's directory and then `findFile`, which scans the
package-level `searchPaths` slice in order. That slice is process-wide and
purely additive: every directory any module has been loaded from stays in it
for the life of the process. Resolution is therefore "first match across
every directory seen so far," and it depends on import *order* — two files
sharing a basename in different folders resolve to whichever directory was
visited first, which can change when an unrelated import is added
elsewhere. The fault's own help text ("modules are looked for next to the
file importing them") describes the intent, not the implementation.

Compounding it, the two import forms are near-homographs for different
operations: `import "path" as name` binds the whole module, while
`import name from "path"` binds one named export out of it. The safe rule,
used throughout Chisel and Studio, is to always import by full path from the
project root — which works, but is a convention papering over resolution
that should be deterministic on its own.

**Severity: mid.** Suggested fix: resolve relative to the importing file
only, and treat the accumulated global path as the compatibility fallback if
one is needed at all.

### 13.20 Smaller surprises, individually minor

Grouped because each is a one-line note rather than a callout of its own,
and because several are working as intended and want documenting rather than
fixing:

- **`continue` is a keyword**, so it cannot name a method — a tool with
  start/continue/finish verbs has to rename the middle one. Recorded here as
  working-as-intended, which was half right: §13.24 promotes it, because the
  same parser rule also blocks the keyword at every *call site*, on objects
  that have nothing to do with the keyword's meaning.
- **Circular imports are a hard fault.** Genuinely useful (§8.9 reports the
  cycle rather than half-loading a module), but it does dictate the file
  graph of anything large, and that consequence is not written down.
- **`++`/`--` are postfix only.** `++x` is a syntax error. Deliberate, per
  §8.4; undocumented as a restriction.
- **A line beginning with `[` or `(` continues the previous statement.**
  Already recorded as §13.12 bug #1 and accepted there as documented
  behavior; relisted here only because it was rediscovered independently,
  which is evidence the §8.1 wording is not doing its job.
- **`%=` is missing** while `+=`, `-=`, `*=` and `/=` all work — see §12,
  where it is filed as a gap to close rather than a defect.
- **`list.length` without parentheses** now raises a proper `property error`
  rather than silently handing back the method, so the hazard reported from
  the Chisel work no longer exists — noted here because that report is
  otherwise cited as a whole and this one item of it is stale.

### 13.21 `and`/`or` do not short-circuit, and the guard idiom every neighbouring language teaches therefore crashes — done

```ghost
target = null

if (target == null or target.hint == '') {   // was: property error, cannot read
    return null                              // property `hint` of null.
}                                            // now: returns null, as written
```

§8.4 used to state this outright — "no built-in short-circuit special-casing
beyond ordinary infix evaluation order: left is evaluated, then right, then
combined" — and `evaluator/infix.go`'s `evaluateInfix` matched it: both sides
were evaluated, and only then did the switch reach `evaluateBooleanInfix`.
So unlike §13.15 and §13.17 this was never drift; the implementation did what
this document asked for. **The callout was against the stance, not the code**,
which is why closing it took a §14 decision (11) and not just a patch.

The stance is wrong for one reason that no amount of documenting fixes.
Ghost takes `and`/`or` from Python and Ruby, and the rest of a reader's Ghost
looks like JavaScript or PHP. All four of those languages short-circuit
`&&`/`||`, so the null guard every one of them teaches —
`x == null or x.field`, `x != null and x.method()` — reads as correct Ghost,
passes review, and then raises the exact fault it was written to prevent, at
the moment the guarded value is actually null. A bare truthy guard fails
identically: `if (x and x.foo)` dereferences `x.foo` whether or not `x` is
falsy.

This is the only finding in this block that reached a real run twice. Chisel
shipped it in `Ui.paintTooltip()`, and then — the dangerous one — in
`Keymap.dispatch()` as `route == null or !this.passes(route.middleware,
studio)`, which crashed on every keypress that had no binding, meaning very
nearly every keypress. An audit of every `and`/`or` in that codebase, run
once the shape was known, found six more live instances (`Ui.focus`,
`Keymap.bind`, `Keymap.passes`, two command-availability guards,
`Signals.forget`) — none yet triggered, all real. Eight defects in one
medium-sized library, from one operator behaving unlike its spelling.

Two things make the fix cheaper here than the same change would be in Python
or JavaScript:

- **There is no value-returning question to settle.** `and`/`or` in Ghost
  are strictly boolean-to-boolean — `1 and 2` is a type error today, not `2`
  — so short-circuiting means returning `false` from a false left operand
  and `true` from a true left one. Ghost never has to decide whether `a or b`
  yields `b`, which is the part of this semantics that makes JS's `||`
  and Python's `or` subtle.
- **The evaluator already evaluates lazily elsewhere.** `cond ? a : b`
  evaluates one arm, verified against the interpreter, so laziness is not a
  new concept in this tree-walker — it is one `case` that has not been
  written.

**Fix.** `evaluateInfix` (`evaluator/infix.go`) now routes `token.AND` and
`token.OR` to `evaluateLogicalInfix` (`evaluator/boolean.go`) as soon as the
left operand is evaluated, before the right one is touched. That function
requires the left operand to be a `Boolean`, returns immediately when it
settles the answer (`false and x`, `true or x`), and otherwise evaluates the
right operand and answers with it — because once the left operand has not
decided, the result *is* the right one: `true and x` is `x`, `false or x` is
`x`. `and`/`or` are gone from `evaluateBooleanInfix`, which is handed both
operands already evaluated and so cannot make this decision; the cases there
would be dead code. `optimizer/fold.go`'s `foldBooleanInfix` keeps folding
two literal booleans — with no evaluation to skip either way — and only its
comment changed.

**Errors name the side at fault.** A non-boolean operand cannot be reported
as `cannot use `and` between null and boolean` any more, because when the
left operand is wrong the right one has deliberately not been evaluated and
naming a type it might have had would be inventing one. `logicalOperandError`
(`evaluator/errors.go`) reports `cannot use `and` with null on the left`
instead, on either side, and adds `help: compare it first, as in `x != null``
when the offending operand is null — which is nearly always the truthy-guard
mistake this whole callout is about. The gain is not only wording: a bare
`if (x and x.foo)` now fails at the `and` with a `Type` fault, where it
previously died with a `property error` on the very dereference the guard
existed to prevent.

**One behavior loosens**, as §14 decision 11 said it would: an unreached
operand is no longer type-checked, so `false and 1` is now `false` where it
was a type error. That is the same trade every short-circuiting language
makes, and it is the point — the unreached side is unreached.

**Severity: high, shipped a crash twice.** Tested in
`evaluator/logical_test.go`: the truth table (unchanged), short-circuiting
proved by a right operand that raises and by one whose side effect is counted
and must not happen, the reached/unreached type-error wording with positions,
and the null help line. All 41 programs in `examples/` produce identical
output before and after — `mud.gs` differs only in how many frames its
non-terminating interactive loop renders inside the timeout, with its output
byte-identical up to that point.

### 13.22 A method's name shadows a same-named import, in every method of its class

```ghost
import "ghost:math" as math

class Theme {
    load() { return math.floor(3.7) }   // property error: function has no
                                        // method `floor`
    math(role) { return role }
}
```

A method body's scope is the class environment (§8.8), and resolution
reaches that environment before it reaches the file's top-level scope where
imports are bound — the same mechanism that lets one method call a sibling
by bare name. A method named `math` therefore shadows `import "ghost:math"`
for *every* method in the class, not only the one that declares it, and the
name resolves to the method object rather than the module.

The report is the worst part. It is a `property error` naming the member
that was called (`function has no method \`floor\``), pointing at the call
site, with nothing at the import and nothing at the method declaration to
connect the two. The reader is looking at a correct call to a module they
correctly imported.

This shipped in Chisel's `chisel/theme.gs`, where `Theme.font(role)`
shadowed `import "lumen:font"` inside `Theme.loadFonts()`, and it survived
that repository's test suite because the collision needs an engine module to
reproduce and `ghost test.gs` has no engine. It surfaced on the first real
run. The workaround is an alias chosen not to collide
(`import "lumen:font" as fontModule`), which requires knowing the hazard
exists.

**Severity: high, shipped a crash.** Suggested fix: a diagnostic, not a
change to resolution order — the resolution order is what makes bare sibling
calls work and should stay. Reporting the collision where it is created, as
§13.18 proposes for a field and a method sharing a name, catches both of
these class-shadowing hazards with one check at class construction.

### 13.23 A function held in a field cannot be called through the field

```ghost
class Command {
    constructor() { this.guard = function () { return true } }
    run() { return this.guard() }   // property error: class `Command` has
}                                   // no method `guard`
```

`x.field(...)` parses as a method call, and method lookup consults the class
chain, which never sees instance fields — the mirror image of §13.18, where
a property *read* consults the instance and so never sees the method. The
workaround is to bind the field to a local first (`test = this.guard;
test(studio)`), which works and reads like an apology.

Callbacks held as fields are an ordinary shape — a command's guard, a
validator, a comparator, a widget's `on('click')` handler — and this makes
every one of them awkward at its call site. Caught by Chisel's own tests
failing on `Command.isEnabled()`, so it is loud rather than silent, which is
the only reason it ranks below §13.22.

**Severity: mid, loud.** Suggested fix: when method lookup fails, consult the
instance's fields for a callable before raising, so `this.guard()` falls back
to calling the field. §13.18's diagnostic keeps that unambiguous by rejecting
the case where both exist.

### 13.24 A reserved word is unusable as a method name *and* at any call site

§13.20 records that `continue` cannot name a method. The sharper half was
missed there: because the parser rejects the keyword after a `.` as well, a
reserved word cannot be *called* on an unrelated object either.

```ghost
cursors.use('arrow')   // syntax error: expected a name, found `(`
```

`use` is only meaningful inside a class body, where it pulls in a trait, yet
it blocks the name everywhere — including on an object that has nothing to
do with traits and that a third-party library may already have shipped.
Chisel's tools have a `drag` verb rather than a `continue` one for exactly
this reason.

The error compounds it: `expected a name, found \`(\`` points at the
parenthesis, one token past the word that actually caused the failure, so it
reads as a malformed call rather than a reserved name.

**Severity: mid.** Suggested fix: accept a keyword as the member name after a
`.` — it is unambiguous in that position — which removes the call-site half
entirely. Whether a keyword may *declare* a method is a separate and smaller
question. Failing that, the fault should at least point at the keyword and
say which one it is.
---

## 14. Decisions for 1.0

The following are this specification's working answers to the
product-direction calls §12 and §13 raise but can't resolve on their own —
made here so 1.0 ships with an answer rather than an asterisk. Revisit any of
them only if real usage argues otherwise; until then, this is what 1.0
targets.

1. **User-defined function arity — done, revised.** Functions and methods
   defined in Ghost get the same missing-argument checking every library
   function already has (§8.7, §12): a call that leaves a required
   parameter unbound becomes an `Argument` fault naming the call, exactly as
   it already does for a library method — closing the largest remaining
   behavioral gap between user code and library code, and the one most in
   tension with the "no silent gaps" goal in §3.

   The original version of this decision also rejected calls with *too
   many* arguments. That half was reverted: it put user-defined functions
   and `object.Function.Evaluate` (the path `list.map`/`filter`/`reduce`/
   `each`/`sort` call a callback through) at odds with each other -
   `Evaluate` never enforced an upper bound, so a callback like
   `(item) => item * 2` already worked as `list.map`'s argument even though
   `map` also passes an index, only because that path skipped arity
   checking entirely. Requiring every trailing parameter a callback doesn't
   use to be declared anyway (`(item, index, list) => item * 2`) fought that
   existing convention rather than matching it, for no benefit worth the
   friction. The maximum is gone for every user-defined function now, not
   only callbacks: extra arguments are dropped, the same way `Evaluate`
   already dropped them. The minimum is unchanged and still enforced.
   Implemented in `evaluator/function.go` (`checkArity`,
   `createFunctionEnvironment`), tested in `evaluator/evaluator_test.go`'s
   `TestFunctionArity` and `TestFunctionArityAllowsDefaultsAndVariety`.
2. **`Map` iteration order — done.** `Map` guarantees insertion order for
   `keys()`, `values()`, `entries()`, `for ... in`, and `String()` (§13.5),
   matching the predictability JS objects and PHP associative arrays already
   give their users (§2). Implemented by backing `Map` with an
   insertion-ordered structure (`object/map.go`'s `order` field and
   `SetPair`/`RemovePair`/`OrderedPairs`) rather than a bare Go map, and
   fixing `ast.Map.Pairs` to preserve source order too, since a map
   literal's order needs to survive parsing before evaluation can preserve
   it.
3. **`string.find`/`findAll`/`matches` — done.** These now take the
   conventional `subject.find(pattern)` shape (§13.7), matching JS, PHP, and
   Python, and the "expressive, familiar syntax" goal in §1. `findAll` was
   fixed at the same time to return every match, not just the first.
4. **`http`'s scope.** `http` stays intentionally minimal for 1.0 — a
   single handler/listen pair, good enough to demonstrate embedding a
   server, not a web framework. Response objects, headers, status codes,
   cookies, and routing are explicitly out of scope for 1.0 and can be
   reconsidered for a later release if real usage asks for them.
5. **OOP surface.** No static members, access modifiers, interfaces, or
   abstract classes for 1.0 (§4, §8.8). Ghost's class model stays
   deliberately small — every member public, one level of inheritance,
   traits for the rest — as a permanent design stance for the 1.x line, not
   a placeholder.

   Reaffirmed, with the cost now measured rather than assumed. The Chisel
   and Studio work (§13.13–§13.20) reports the absence of statics as the
   single thing that most distorts the shape of a library built on Ghost:
   with no named constructors, `Rect.fromBounds(...)` becomes a free
   function in a helper file, and with no class constants, `Button.HEIGHT`
   becomes a key on a theme object. Both workarounds are real and both were
   used throughout. Neither is unworkable, and neither argues that Ghost's
   class model is wrong — they argue that a library's factories and
   constants have to live somewhere less obvious than the class they belong
   to. That is a genuine cost, recorded here so the stance is upheld with
   its price written down rather than by not having looked. Revisit for 1.x
   only if a second independent library reports the same distortion.
6. **Cross-type `==`.** Cross-type `==`/`!=` keeps raising a type error
   rather than returning `false` (§8.5). It is consistent with "operators
   keep one meaning" and already covered by a test; what this decision
   obligates is documentation, not implementation — the behavior needs a
   prominent, early callout in user-facing docs (a getting-started page, not
   only this specification) so it reads as a documented design choice
   rather than a surprise.
7. **Multi-instance embedding.** Hosting more than one independent Ghost
   configuration in a single Go process is not a 1.0 requirement.
   `RegisterFunction`/`RegisterModule` stay process-global (§10.3); a host
   that needs isolated configurations runs separate processes for 1.0.
   Revisit only if an embedding use case actually needs it.
8. **Date time zones.** This reverses the earlier UTC-only design (previously
   documented directly in §9.5 and `object.Date`'s doc comment, not merely
   proposed): a `Date` now carries an explicit time zone, defaulting to UTC,
   moved with `date.inTimeZone()` or set at construction with
   `date.ofInZone()` (§9.5). The reproducibility guarantee that motivated the
   original design — a date built once means the same thing no matter which
   machine runs the script — is kept in full for the operations it actually
   protects: instant comparison (`< <= > >= == !=`) and every function that
   reduces to comparing or measuring instants (`differenceInX`, `toUnix`,
   `isSameDay`, ...) stay zone-independent, unaffected by this change. What
   changes is that a *calendar* reading (`year()`, `format()`, `toString()`,
   `startOfDay()`, ...) can now be asked for in an explicit zone. Reproducing
   the original problem would require reading the host machine's *ambient*
   configured zone with no name in the script naming it — that stays absent
   by design (no `date.local()` or equivalent); every zone this module
   resolves is a name the script itself wrote down, looked up against the
   IANA tz database Ghost embeds in its own binary (`time/tzdata`), so the
   same name resolves the same way on every platform Ghost ships for. A
   script that never calls `inTimeZone`/`ofInZone` behaves exactly as it did
   under the old design.
9. **The scoping model — decided, and done.** §13.13, §13.14 and §13.15 were
   one question asked three ways: *when a name is assigned, and the name
   already exists further out, what happens?* The answer, chosen over adding a
   declaration keyword, is:

   **Assignment walks the scope chain, and blocks have scopes of their own.**
   `x = 5` rebinds the nearest existing `x` in the enclosing chain and
   declares a new one only when the name is bound nowhere; each block —
   `if`/`else`, `while`, `for`, `for ... in`, a matched `switch` case —
   introduces an environment, and a loop binds its control variable once per
   iteration. §8.3 documents the result.

   The two halves are one decision because each pays for the other. Walking
   assignment alone fixes §13.13 but has a mirror-image hazard: with no way to
   say "this one is mine," every function-local temporary becomes a potential
   write to an outer name that happens to match, and inside a method the chain
   reaches the class environment, so a local named `scale` would rebind a
   *sibling method* named `scale`. Block scoping is what contains that — a
   name first assigned inside a block stays inside it — and it is what §8.3
   promised all along. Block scoping alone, meanwhile, would have made
   §13.13's silent failure worse rather than better, since more scopes means
   more places an assignment can fail to reach.

   The rejected alternative was a declaration keyword (`let`, or `var`):
   assignment stays local, and an explicit form makes shadowing deliberate.
   Safer semantics, and the one every reader of modern JS already has loaded.
   It was not chosen because it costs a reversal of §8.3's most prominently
   documented stance, a new keyword in a language that has kept its surface
   deliberately small (§4), and a migration for every line of existing Ghost
   including this document's own examples — to buy protection against a
   shadowing bug that block scoping already contains.

   What this decision cost, recorded honestly: one breaking change (§13.15 —
   a name first assigned in a branch no longer outlives it, though all 41
   programs in `examples/` are byte-identical before and after), the sibling-
   method rebinding hazard above, and 1.03×–1.21× wall time from the extra
   environment link every name lookup crosses. Allocation is at parity, which
   took the scope reuse described in §13.15.

10. **Module export surface — open (§13.16).** Related to decision 9 in
    spirit but separable in practice: whether Ghost gains an explicit
    `export`, or a naming convention, or keeps exporting everything. The
    fallback that keeps every existing file working is to treat a module
    that marks nothing as exporting everything, exactly as today.

11. **`and`/`or` short-circuit — reversing §8.4 (§13.21) — done.** §8.4's
    non-short-circuiting rule was a decision this document made and the
    interpreter honoured; §13.21 is the evidence against it, and the
    evidence is strong enough to reverse it for 1.0. Eight real defects in
    the first library written against Ghost, two of which shipped, all with
    the same shape: a null guard written the way Python, Ruby, JavaScript
    and PHP all teach it, dereferencing the value it just tested.

    The argument for the current rule is "operators keep one meaning" — the
    same principle that keeps `+` from concatenating a number onto a string
    (§"Error handling") and keeps `<` from ordering two lists (§8.4). It
    does not apply here. Short-circuiting does not give `and` a second
    meaning; it gives it the *same* meaning, computed without evaluating an
    operand whose value cannot change the answer. `false and x` is `false`
    for every `x`, which is exactly why the right operand can be skipped.

    Because Ghost's `and`/`or` are boolean-only, the reversal is narrow:
    they keep answering a `boolean`, they keep raising a `Type` fault on a
    non-boolean operand *that is reached*, and the only observable loosening
    is that an unreached operand is no longer type-checked. Ghost does not
    inherit the value-returning semantics that make this operator subtle
    elsewhere, and should not add them.

    The rejected alternative was to keep the behavior and document it
    harder. §13.21 is what that costs: the behavior was already documented,
    in this section, in the sentence §13.21 quotes — and it still shipped
    twice, because the failing code reads correctly to anyone who knows any
    of the four languages Ghost is presented as resembling. A rule that
    survives its own documentation is a design defect, not a teaching
    problem.

    Implemented in `evaluator/boolean.go`'s `evaluateLogicalInfix`, reached
    from `evaluator/infix.go`; §8.4 now documents short-circuiting as the
    rule and §13.21 records what changed. The narrow reversal held: the truth
    table is untouched, both operands are still booleans, and the single
    observable loosening is the unreached operand going unchecked. One thing
    the decision did not anticipate is that the error *wording* had to move
    too — a wrong left operand can no longer be reported against a right
    operand that was deliberately never evaluated, so `and`/`or` now name the
    side at fault (`cannot use `and` with null on the left`). That reads
    better than what it replaced: the truthy-guard mistake now fails at the
    operator instead of at the null dereference downstream of it.

---

## 15. Fix Priority for the Chisel/Studio Findings

§13.13–§13.24 arrived as one report — `docs/papercuts.md` in the Studio
repository — rather than as separate findings, so they have never been
ranked against each other. This section does that, and is the entry point
for a session picking up this work: take the highest open item.

The ranking is by expected damage — how badly it fails, times how likely
ordinary code is to hit it — with cost used only to break ties. It is not
the order they were discovered in, and it deliberately promotes §13.21
above findings that have been open longer.

### Closed since the report was written

Five items in `papercuts.md` no longer reproduce against this interpreter.
The first four were verified by running each one at `c31c79d`; §13.21 was
closed here, by this section's own ranking:

| Finding | Papercut severity | State |
|---|---|---|
| §13.13 assignment does not reach an enclosing scope | high, silent | Fixed — §14 decision 9 |
| §13.14 closures cannot capture a loop variable | high, silent | Fixed — §14 decision 9 |
| §13.15 blocks do not introduce a scope | high, spec drift | Fixed — §14 decision 9 |
| `list.length` hands back the method | low | Fixed — §13.20 |
| §13.21 `and`/`or` do not short-circuit | high, shipped twice | Fixed — §14 decision 11 |

That report is cited elsewhere as a whole; these four are stale, and the
architecture notes justifying workarounds for them (state on instances,
`make…` closure factories) now describe a constraint that is gone — though
the workarounds themselves stay correct, so nothing built on them has to
change.

§13.15's breaking change was also checked against the codebase that reported
it rather than only against `examples/`: Studio's engine-independent suite,
132 cases, passes unchanged at `c31c79d`. `Dock.arrange()` — the one place
that report names as depending on block-free scoping — is directly covered
and unaffected, because it binds `area` and `taken` before the `switch`, so
the destructuring assignment inside each case rebinds them through the
walking assignment §14 decision 9 introduced alongside block scoping. The two
halves of that decision paying for each other is not theoretical.

### Open, in priority order

| # | Finding | Severity | Cost | Why here |
|---|---|---|---|---|
| ~~1~~ | ~~§13.21 `and`/`or` do not short-circuit~~ | high | small | **Done** — §14 decision 11. |
| 2 | §13.22 a method's name shadows a same-named import | high | small | Shipped; silent until the call runs; the fault points at correct code. A diagnostic at class construction is the whole fix, and §13.18 wants the same check. |
| 3 | §13.17 a bare sibling call loses the receiver | mid | small | §8.8 actively teaches the broken form, so the document is generating the bug. Contained in `unwrapCall`. |
| 4 | §13.18 a field and a method may share one name | mid | small | Silent, and the behavior is the opposite of what a reader assumes. Shares its fix site with #2 — do them together. |
| 5 | §13.16 every top-level name is exported | high | large | Highest damage left, but it is a language design decision (§14 decision 10 is still open) rather than a defect with a known patch. Blocks nothing today; distorts every file layout built on Ghost. |
| 6 | §13.23 a field-held function cannot be called through the field | mid | small | Loud, and the workaround is one line, which is the only reason it sits below the silent findings. |
| 7 | §13.24 a reserved word is unusable at a call site | mid | mid | Parser change; unblocks names a third-party library may already use. The fault pointing one token late is worth fixing even if the rest is not. |
| 8 | §12 `math.floorDiv(a, b)` | mid | trivial | A table entry. Every line of pixel layout wants it. |
| 9 | §12 `%=` | low | trivial | A table entry; the only compound operator missing. |
| 10 | §13.19 module resolution is global and first-match-wins | mid | mid | Order-dependent and able to change under an unrelated import, but a full-path convention avoids it completely, and no reported bug has come from it yet. |

**The next item is #2.** #1–#4 were four small, independent patches against
known code paths that together close every finding whose failure does not
point at its own cause — #1 reported at the dereference, #2 at the call site,
#3 at the callee, and #4 reports nothing at all. #1 is now done; #2 and #4
should still land as one change, since a single check at class construction
catches both.

### Already answered, no work outstanding

- **No statics** — §14 decision 5, upheld with its cost recorded. The
  Chisel report is the first of the two independent reports that decision
  names as its condition for revisiting; it is not a defect.
- **Division promotes to float** — not reproducible as stated. `6 / 3` is
  `2`, and §8.4's int/float rules work as documented; what is missing is a
  way to *ask* for the integer quotient, filed as #8 above.
- **Circular imports are a hard fault**, **`++` is postfix only**, **no
  `const`**, **a line opening with `[` or `(` continues the previous
  statement** — all working as intended (§13.20, §13.12). They want a
  user-facing mention in the getting-started documentation, which §14
  decision 6 already owes for cross-type `==`, rather than a code change.
