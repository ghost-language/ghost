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
product-direction calls needed to close both gaps.

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
- **No script can crash the host.** This is already true in the general case
  (`ghost.Execute` recovers panics — see §8.11) but the RNG data race in §13.1
  is a live counterexample under concurrent use, and closing it is a
  precondition for the tag.
- **The CLI and its own documentation agree with each other and with the
  code.** They currently do not (§13.4, §13.10).
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
  a `.ghost` file on disk next to the importing file, or a standard library
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
breaking change for any `.ghost` script written against the old name, made
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
  present; otherwise a statement simply ends where its grammar says it ends
  (there is no significant-newline rule and no automatic-semicolon-insertion
  logic to reason about).

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

There is no declaration keyword — `x = 5` both declares and assigns. Scoping
is lexical: closures capture their defining environment, and blocks
(`if`/`while`/`for`/function bodies) each introduce a new enclosed
`Environment` chained to its parent. A `for`/`for ... in` loop's control
variable(s) are scoped to the loop and restored (or removed, if they did not
exist before) to whatever they were bound to outside it once the loop ends —
so a loop variable does not leak a stray binding into surrounding code, but
*does* transparently shadow-and-restore an existing variable of the same name.

There is no destructuring assignment (`[a, b] = list` or `{a, b} = map` are
not supported — see §12) and no chained assignment (`a = b = 5` does not
parse as one assignment to both).

### 8.4 Operators

| Category | Operators | Notes |
|---|---|---|
| Arithmetic | `+ - * / %` | On numbers: standard, with the int/float promotion rules above. On lists: elementwise with **NumPy-style broadcasting** — see below. On strings: only `+` (concatenation); `-`/`*`/`/`/`%` on strings are a type error. |
| Comparison | `< <= > >=` | Numbers and strings only (strings compare lexicographically). **Not supported between two lists** — deliberately: neither an elementwise nor a lexicographic reading was judged obviously correct (`CLAUDE.md`). Dates support `< <= > >=` as instant ordering. |
| Equality | `== !=` | See §8.5 — this is one of the language's most distinctive (and most incomplete — §13.2) behaviors. |
| Logical | `and`, `or`, `!` | Word operators, not `&& \|\|` — there is no `&&`/`\|\|` token at all. `!` is the only prefix logical operator. Both operands of `and`/`or` are evaluated as ordinary booleans (no built-in short-circuit special-casing beyond ordinary infix evaluation order: left is evaluated, then right, then combined). |
| Unary | `-`, `!` | `-` negates a number only. `!` follows Ghost's truthiness rules (§8.5), not "must be boolean." |
| Range | `a..b` | Inclusive integer range, producing a `list`: `1..5` → `[1, 2, 3, 4, 5]`. Descending (`a > b`) produces an empty list rather than counting down. Not foldable at compile time (would require a shared mutable literal). |
| Ternary | `cond ? a : b` | Standard. |
| Assignment | `=` | Also declares. Valid targets: a bare identifier, an index expression (`list[0] =`, `map["k"] =`), or a property expression (`instance.field =`, `map.key =`). |
| Compound assignment | `+= -= *= /=` | **No `%=`.** Desugars to `target = target OP value`. |
| Increment/decrement | `++ --` | Postfix only (`x++`, not `++x`); operates only on a variable holding a number. |
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

**Equality (`==`/`!=`) is comparison, not coercion**, and behaves differently
depending on the pair of types involved:

| Left / Right | Result |
|---|---|
| Same primitive type (number, string, boolean) | Value equality, as expected. |
| Either side `null` | `true` only if *both* sides are `null`; otherwise `false`. This is the one cross-type comparison that is allowed and does not error. |
| Both `list` | **Deep structural equality**, to any depth (`object.ValuesEqual`/`ListsEqual`). |
| Both `instance` | **Identity** — two separate instances with identical fields are not `==`. |
| Both `date` | Instant equality. |
| Both any other same type (`map`, `function`, `class`, `trait`, ...) | **Type error** — see §13.2; this must be fixed before 1.0, not a design choice. |
| Different, non-null types (`5 == "5"`, `[1] == {}`) | **Type error**, not `false`. This is deliberate and covered by an explicit test (`evaluator_test.go`), consistent with `CLAUDE.md`'s "operators keep one meaning" principle. §14 confirms this stays for 1.0 — see that section for what it obligates the documentation to do. |

### 8.6 Control Flow

- **`if (cond) { } else if (cond) { } else { }`** — parentheses around the
  condition and braces around every branch are both mandatory; there is no
  brace-less single-statement form.
- **`while (cond) { }`**.
- **`for (i = 0; i < n; i++) { }`** — C-style three-clause form. The increment
  clause accepts an assignment, a compound assignment, or a postfix
  increment/decrement, but not an arbitrary expression.
- **`for (key, value in iterable) { }`** and **`for (value in iterable) { }`**
  — iterates a `list` (key = integer index) or a `map` (key = the map key, in
  insertion order once §13.5 is fixed — today it is Go's randomized map
  order). Iterating anything else is a type error with the help text
  *"`for ... in` walks a list or a map."*
- **`switch (value) { case a { } case b, c { } default { } }`** — this is a
  **match-expression**, not a C-style switch: there is no fallthrough, no
  `break` is needed or accepted between cases, and a `case` may list several
  comma-separated values that all run the same block. At most one `default`
  is allowed (a second is a parse error). See §13.3 for a correctness gap in
  how case values are compared.
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
```

- Functions are first-class values, close over their defining scope, and may
  be named (hoisted into the enclosing scope/class at definition time) or
  anonymous.
- Parameters support **default values** (`name = "Hello"`), evaluated lazily
  per-call in the function's own scope (so a default can reference an
  earlier parameter or an enclosing binding). Defaults are *not* required to
  trail all non-default parameters — the parser does not enforce an
  ordering — worth deciding whether that should be a parse-time error for
  1.0.
- **No rest/variadic parameters** (`...args`) and **no spread** in a call or
  a list literal (`f(...args)`, `[...list]` are not supported). A caller who
  passes too many arguments simply has the extras ignored; too few leaves
  the missing parameters unbound (reading them raises the ordinary
  "not defined" name error, not a dedicated arity error) — §14 decides that
  user-defined functions get the same strict arity checking every library
  function already has, closing that gap (§12).
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
  each other by bare name** with no `this.` required (though `this.method()`
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
import "helpers"                       // a .ghost file, whole-file import for side effects
import add, subtract from "math_ext"   // named imports from a .ghost file
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

import "lumen:image", { Spritesheet }  // binds `image` *and* `Spritesheet` in one statement
image.something()                      // the whole module, same as a bare `import "lumen:image"`
new Spritesheet("sheet.png")           // the named export, same as `import { Spritesheet } from ...`
```

**File imports** (no scheme — any other string) name a `.ghost` file:

- A module is a `.ghost` file, looked up **next to the file that imports
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
- The combined form works here too: `import "helpers", { greet }` binds
  `helpers` (a `Map` of the module's top-level bindings, same as the bare
  form) and `greet` in one statement.

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
  named import from a `.ghost` file. A property (`pi`) is evaluated once,
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
- **The two forms combine**: `import "path", { a, b }` binds the whole
  module — exactly as the bare form does, alias and all (`import "path" as
  m, { a, b }` also works) — *and* pulls `a`/`b` out of it by name, in one
  statement. This is the fix for needing both the module itself and one of
  its members: `import "lumen:image", { Spritesheet }` binds `image` (so
  `image.something()` works, if the module has any methods of its own) and
  `Spritesheet` (so `new Spritesheet(...)` does too), where before this took
  two separate `import` lines naming the same path. `import "path", { * }`
  combines the whole module with every export. The trailing name list has
  to be braced here — unlike the standalone `from` form, there is no
  trailing `from` left to mark where an unbraced list would end — and, since
  a standalone function has nothing to pull a name out of, is rejected the
  same way `import { x } from "scheme:someFunction"` already is.
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
 --> example.ghost:4:15
  |
4 | total = count + label
  |               ^
  |
  = in sum(), called at example.ghost:9:1
  = help: both sides of `+` have to be the same type
```

**Nothing below `ghost.Execute` is allowed to reach the caller as a Go
panic.** Runaway recursion is counted and reported as an ordinary value
error before the Go stack actually overflows (§8.7); a bug in Ghost itself
that *does* panic is recovered at the top of `Execute` and reported as an
`Internal` fault asking the reader to file a bug, with the Go stack trace
attached only when `GHOST_DEBUG` is set in the environment. The one
documented exception to "nothing crashes the host" is the RNG data race
described in §13.1, which is a real gap in that guarantee under concurrent
use, and closing it is one of the preconditions for the tag (§3).

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
line to output. `token.PRINT` (§13.9) was already dead code before this
removal and stays exactly as dead after it — that defect is about an unused
token constant, not about `print` the function.

### 9.2 Built-in Methods on Core Types

These are called directly on a value (`"hi".toUpperCase()`), independent of
the module system.

**`boolean`** — `toString()`.

**`null`** — `toString()` (→ `"null"`).

**`number`**

| Method | Arity | Behavior |
|---|---|---|
| `round([places])` | 0–1 | Rounds to the nearest integer, or to `places` decimal places. An already-integral `Number` is returned unchanged (not converted through a float round-trip). |
| `floor()` | 0 | Rounds down. Integral input returned unchanged. |
| `toString()` | 0 | |

*No `ceil()` instance method* despite `round`/`floor` both existing — only
reachable via `math.ceil()` (§12). No `abs()`, `pow()`, `sqrt()`,
`clamp()`, or any of the `isX` predicates as instance methods either; all of
those live only in the `math` module.

**`string`**

| Method | Arity | Behavior |
|---|---|---|
| `find(subject)` | 1 | **The string itself is compiled as a regular expression**; `subject` is the text searched. Returns the first match, or `""` if none. See §13.7 — this shape is scheduled to flip before 1.0. |
| `findAll(subject)` | 1 | As `find`, but returns a `list` of the submatches of the *first* match (despite the name, this does not find every match in the subject — a genuine bug, not a naming choice; it calls `FindStringSubmatch`, not `FindAllString`). |
| `matches(subject)` | 1 | As `find`, returning a `boolean`. |
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

*No* `contains`/`includes`, `indexOf`, `repeat`, `padStart`/`padEnd`,
`charAt`/`at`, `slice`/`substring`, `reverse`, or `isEmpty` — see §12.

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

*No* `indexOf`/`find`/`findIndex`, `flatten`, `some`/`any`, `every`/`all`,
`splice`/`insertAt`/`removeAt`, `fill`, `chunk`, `flatMap`, or `isEmpty` —
see §12. There is also no `push`-equivalent that inserts at an arbitrary
index or a front-insert (`unshift`).

**`map`**

| Method | Arity | Behavior |
|---|---|---|
| `get(key[, default])` | 1–2 | `default` if given and the key is absent, else `null`. |
| `has(key)` | 1 | True if the key is present, regardless of its value (distinguishes "absent" from "present and `null`"). |
| `set(key, value)` | 2 | Mutates in place; returns the map itself (chainable). |
| `keys()` / `values()` | 0 | Returns a `list`, in insertion order once §13.5 is fixed (today, Go's randomized map-iteration order). |
| `merge(other)` | 1 | New map; on a key collision, `other`'s value wins (same rule a later assignment to the same key would follow). |
| `length()` | 0 | |

*No `delete`/`remove` method* — once a key is set, there is no way to
remove it from a Ghost script (§12, and see §13.2 for the related equality
gap). No `entries()`, no `forEach` (use `for ... in`).

**`date`** — `toString()` only (ISO-8601/RFC3339, always UTC). Every other
date operation is a function in the `date` module (§9.5), not a method,
which keeps `Date` itself immutable and side-effect-free and is a
deliberate design choice modeled on `date-fns` (§2) rather than a mutable
built-in `Date` class — see the doc comment on `object.Date`. Dates support
`< <= > >= == !=` directly as operators (instant comparison) but no
arithmetic operators (`date1 + date2` is a type error, with help text
pointing at the `date` module).

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
see §13.1, the most serious defect in this specification, for the
concurrency gap in this shared state).

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
and returns a new value; nothing mutates. Every `Date` is UTC-normalized —
Ghost does not model time zones at all, so a date built once compares the
same everywhere the program runs (the same reproducibility guarantee a
seeded `random` run gets).

**Construction/conversion** — `now()`, `today()` (midnight UTC),
`of(year, month, day[, hour, minute, second])` (month is 1–12; an
out-of-range day/hour/minute/second is a `Value` error rather than silently
rolling into the next period), `parseISO(text)` (RFC3339 or bare
`YYYY-MM-DD`), `fromUnix(seconds)`, `toUnix(date)`, `toUnixNano(date)`,
`format(date, pattern)` (date-fns-style pattern letters, not Go's reference
layout — see below).

**Arithmetic** — `addDays`/`subDays`, `addWeeks`/`subWeeks`,
`addMonths`/`subMonths` (clamps to the target month's last day rather than
rolling over — Jan 31 + 1 month = Feb 28/29, not Mar 2/3),
`addYears`/`subYears`, `addHours`, `addMinutes`, `addSeconds` (all arity 2:
date, count). *No `subHours`/`subMinutes`/`subSeconds`* — asymmetric with
every other pair (§12).

**Predicates** — `isSameDay(a, b)`, `isWeekend(date)`, `isLeapYear(date)`.

**Differences** — `differenceInDays`/`Hours`/`Minutes`/`Seconds(a, b)`,
truncated toward zero (so `differenceInDays(a, b) == -differenceInDays(b,
a)` always, never off-by-one from truncation direction).

**Period boundaries** — `startOfDay`, `endOfDay`, `startOfMonth`,
`endOfMonth`. *No `startOfWeek`/`endOfWeek`/`startOfYear`/`endOfYear`* —
gap relative to the `date-fns` surface this module is explicitly modeled on
(§12).

**Components** — `year`, `month`, `day`, `hour`, `minute`, `second`,
`weekday` (0 = Sunday, matching `date-fns`'s `getDay`).

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
reproducibility before a script ever runs. **This generator is shared with
`math.randomInt`/`math.randomSeed`, and neither module synchronizes access
to it** — see §13.1, the most serious correctness gap in this specification.

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
| `extname(path)` | 1 | The extension, dot included (`".ghost"`), or `""` if there is none. |
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
| `-t` | Prints how long execution took. | Works, but **undocumented** in `helpCommand()`'s own printed help (§13.10). |
| `-i` | *(documented)* Runs a file, then drops into a REPL with the script's environment intact. | **Not implemented** — not registered as a flag at all in `cmd/ghost.go`. Documented in both `README.md` and the CLI's own `-h` output. See §13.4. |

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
declare a class in `.ghost` source (§8.8), and nothing here adds a second.
A bundled Ghost-source "prelude" — classes written in `.ghost`, evaluated
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
pure-Ghost plugin/extension mechanism (e.g., a convention for a `.ghost`
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
- [x] Lists and maps, with broadcasting arithmetic on lists
- [x] `if`/`else if`/`else`, `while`, C-style `for`, `for ... in`
- [x] `switch`/`case`/`default` as a match-expression (no fallthrough)
- [x] `break`/`continue`/`return`
- [x] First-class functions, closures, default parameters
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
- [x] `-h`/`-v`/`-t` CLI flags
- [ ] `-i` CLI flag — documented, not implemented (§13.4)

---

## 12. Required for 1.0: Functionality to Add

The following close the gap between the target this specification sets and
what the implementation does today. Organized by area, in roughly descending
order of how likely each is to be hit by an ordinary user.

**String methods:** `contains`/`includes`, `indexOf`/`lastIndexOf`,
`repeat(n)`, `padStart`/`padEnd`, `charAt`/`at`, `slice`/`substring`,
`reverse`, `isEmpty`. Several of these (`contains`, `slice`) exist on `list`
already, so their absence on `string` is a gap to close, not a decision to
revisit.

**List methods:** `indexOf`/`find`/`findIndex`, `flatten`, `some`/`any`,
`every`/`all`, `splice`/`insertAt`/`removeAt`, `unshift` (front-insert, to
pair with `push`/`pop`/`shift`), `fill`, `chunk`, `flatMap`, `isEmpty`.

**Map methods:** `delete`/`remove` (there is currently no way to remove a
key from a map at all), `entries()` (as a list of `[key, value]` pairs, for
symmetry with `keys()`/`values()`).

**Number methods:** `ceil()` (exists on `math` but not as an instance
method, unlike its sibling `round`/`floor`), and, depending on how far
parity with `math` should go, `abs()`/`pow()`/`clamp()`/the `isX`
predicates as instance methods.

**Date module:** `subHours`/`subMinutes`/`subSeconds` (asymmetric with
every other add/sub pair), `startOfWeek`/`endOfWeek`/`startOfYear`/`endOfYear`.

**Language-level gaps:**
- No destructuring assignment (`[a, b] = list`, `{x, y} = map`).
- No rest/variadic parameters or spread syntax, in either direction (call
  site or list literal).
- No arity checking for **user-defined** functions/methods — every library
  call is strictly arity-checked, but a Ghost-defined function silently
  drops extra arguments and leaves missing parameters as undefined names.
  §14 decides this closes before 1.0.

---

## 13. Required for 1.0: Defects to Fix

Ranked roughly by how much damage each can do, not by how easy it is to fix.
Every finding below was confirmed by reading the relevant code paths
directly (file/line references given); none are speculative.

### 13.1 The shared random-number generator is not safe for concurrent use

`library/modules/random.go` keeps the PRNG behind two unsynchronized package
variables, `seed` and `randomizer` (`*rand.Rand`), read and written directly
by `randomRandom`, `randomSeed`, the exported `SeedRandom`, and — from
`math.go` — `mathRandomInt`/`mathRandomSeed`. None of these hold a lock.

Go's `*rand.Rand` is explicitly **not** safe for concurrent use by multiple
goroutines (only the package-level top-level functions, which use their own
internal mutex, are). Ghost's own `http.handle` (§9.11) runs every request's
callback on its own goroutine, and the codebase elsewhere (`evaluator/import.go`'s
`moduleState sync.Mutex`, `object.Scope.Depth`'s doc comment) shows the team
is already deliberately careful about exactly this class of bug — this one
module appears to have been missed. Two concurrent requests both calling
`random.random()` or `math.randomInt()` race on `randomizer`'s internal
state: under `go test -race` this is a reported data race, and in
production it can silently corrupt the sequence or, rarer, panic inside
`math/rand`. This is the one place in the codebase that can violate
§8.11's "nothing crashes the host" guarantee today. **Fix:** guard `seed`/
`randomizer` with a mutex (mirroring `moduleState` in `import.go`), or swap
to a `*rand.Rand` built over a locked source.

### 13.2 `==`/`!=` cannot compare two maps, functions, or several other same-typed pairs — not even by identity

`evaluator/infix.go`'s `evaluateEquality` special-cases only `NULL`,
`INSTANCE` (identity), and `LIST` (deep). Everything else — same-typed
`MAP`, `FUNCTION`, `CLASS`, `TRAIT`, `SUPER`, `SCOPE`, the `LIBRARY_*`
wrapper types — falls through to `evaluateInfix`'s final `operatorError`,
producing `` cannot use `==` between two maps `` even for `m == m` (the
identical object compared to itself). This directly contradicts
`object/equality.go`'s own doc comment on `ValuesEqual` — *"it is what `==`
means between two Ghost values... a value counted as equal inside a list is
equal everywhere else in the language too"* — which is false in practice:
`object.ValuesEqual` (used by `list.contains()`/`list.unique()`) *does*
fall back to identity comparison for maps and functions, but `evaluateInfix`
never calls `object.ValuesEqual` at all — it has its own, narrower,
hand-rolled `evaluateEquality`. The two equality implementations have
silently diverged. **Fix:** either route `evaluateEquality`'s default case
through `object.ValuesEqual` directly, or explicitly special-case `MAP`/
`FUNCTION`/etc. the way `INSTANCE` already is. Given `map` is one of the two
core collection types, this is the highest-value correctness fix required
for 1.0.

### 13.3 `switch` silently swallows an error in its subject or case expressions

`evaluator/switch.go`'s `evaluateSwitch` calls `Evaluate(node.Value, scope)`
and, for each case, `Evaluate(val, scope)`, **without ever checking
`isError()`** on either result. If the switch's subject (or a case value)
fails to evaluate, the resulting `*object.Error` is compared via `.Type()`/
`.String()` against the other branches instead of being propagated — in the
best case this falls through to `default` (or returns `nil` if there is
none), silently discarding a real error instead of surfacing it the way
every other construct in the language does. Separately, case-value
comparison uses `obj.Type() == out.Type() && obj.String() == out.String()`
rather than `object.ValuesEqual` — a comparison by *string representation*
rather than by value, which is both weaker than the rest of the language's
equality rules (§8.5) and unable to correctly match composite values.
**Fix:** check `isError()` immediately after evaluating the subject and
each case value and propagate; route the comparison through
`object.ValuesEqual`.

### 13.4 `-i` is documented in two places and implemented in none

Both `README.md` ("Interactive mode... pass the `-i` flag") and
`cmd/help.go`'s own `-h` output ("`-i` enter interactive mode after
executing file") describe a flag that `cmd/ghost.go` never registers
(`flag.BoolVar` is called for `flagHelp`, `flagVersion`, `flagTime` only)
and there is no code path anywhere that runs a file and then starts a REPL
with its environment intact. A user following either piece of Ghost's own
documentation for this feature will find it silently does nothing (the flag
is simply unrecognized... actually, with Go's `flag` package, an
unregistered `-i` will make `flag.Parse()` itself error and exit — worth
verifying that failure mode is acceptable, since it is not the same as "the
flag is ignored"). **Fix:** implement `-i` before 1.0 ships, since it is
already committed to in the documentation users will read first.

### 13.5 Map/list `for ... in` and `keys()`/`values()` iterate in random order

`evaluator/for_in.go`'s map branch and `object.Map`'s `keys()`/`values()`/
`String()` all range directly over the underlying Go `map[MapKey]MapPair`,
whose iteration order Go deliberately randomizes per-run. Two runs of the
identical script over the identical map can print keys/values/`String()` in
different orders. This is inconsistent with the rest of the language's
posture toward predictability (`list.sort()` is explicit and stable; the
"did you mean" suggestions are tie-broken deterministically) and is a
common source of "flaky" script output and non-reproducible test
assertions for users. §14 decides Ghost guarantees insertion order for 1.0.
**Fix (larger effort):** back `Map` with an insertion-ordered structure (as
JS objects and PHP associative arrays both behave — squarely in Ghost's
stated inspirations, §2) — or, at minimum, sort `keys()`/`values()`/
`String()`'s output deterministically even if internal storage stays a Go
map.

### 13.6 Bounds-checking is inconsistent across list/map/string operations

`list[i]`/`map[k]`/`string[i]` indexing all return `null` for an
out-of-range index or missing key (`evaluator/index.go`), as do
`list.pop()`/`shift()`/`first()`/`last()`/`tail()` on an empty list — but
`list.slice()` raises an `Index` error for an out-of-range `start`/`end`
(`object/list.go`). Neither behavior is wrong on its own, but the split is
not documented as an intentional rule anywhere, which makes it easy for a
new list/string method to pick the "wrong" one of the two conventions by
accident. **Fix:** state the rule explicitly (e.g., "a read that names a
*position* is lenient; an operation that names a *range* validates it") in
§7, and audit existing methods against it before 1.0.

### 13.7 `string.find`/`findAll`/`matches` invert the usual receiver/argument relationship

`(pattern string).find(subject string)` compiles the **receiver** as the
regular expression and matches it against the **argument**. Every widely
used language Ghost is modeling itself on does the opposite —
JS's `subject.match(pattern)`, PHP's `preg_match(pattern, subject)` (subject
still second, but the *call* reads pattern-first, not method-chained on the
subject), Python's `re.search(pattern, subject)`. A Ghost user writing
`text.find(pattern)` out of that muscle memory gets a working call that
does something different from what they meant, with no error to catch the
mistake (both operands are strings). This is called out and clearly
intentional in the code's own comment ("the string a pattern method is
called on is the pattern itself"), but §14 decides it flips to
`subject.find(pattern)` before 1.0 locks the API, while it is still a
breaking change nobody has shipped against yet.

Separately: `findAll` is misleadingly named — it calls
`FindStringSubmatch` (the *first* match's capture groups), not
`FindAllString`/`FindAllStringSubmatch`. A user reaching for "all the
matches in this text" via `findAll` will get only the first match's
submatches. This is a genuine bug, fixed at the same time as the
receiver/argument flip above.

### 13.8 `ghost.extend()` is unavailable on Windows

Go's `plugin` package (used by `ghost.extend`, §9.12) only supports
Linux and macOS; there is no Windows implementation at all (it's a build
failure/`ErrNotSupported` situation, not a degraded one). Ghost itself
builds and ships Windows binaries (`Makefile`, `.goreleaser.yml`). The only
in-language way to extend Ghost with native code is therefore unavailable
on a third of Ghost's own supported platforms, with no fallback and no
documented caveat. At minimum, this needs a documented caveat before 1.0;
a real fallback is a larger question outside this specification's scope.

### 13.9 Dead token: `token.PRINT`

`token/token.go` defines `PRINT` as a token type, but the scanner's
`keywords` map (`scanner/scanner.go`) never maps the string `"print"` to it
— so this token can never actually be produced by the scanner, and nothing
in the parser or evaluator references it either. Ghost has no `print` at all
in 1.0 (§9.1) — `console.log(...)` replaced it — so there is no live
identifier this token could even be mistaken for. `PRINT` is pure dead code,
presumably left over from an earlier design. **Fix:** delete it (and its
`typeNames` entry).

### 13.10 `-t` is implemented but undocumented

The reverse of §13.4: `cmd/ghost.go` registers and honors `-t` ("display
how long the program ran for"), but `cmd/help.go`'s `helpCommand()` never
mentions it in its printed usage/flags list. **Fix:** add it to the printed
help before 1.0.

### 13.11 Minor: duplicated `isTruthy` logic

`evaluator/evaluator.go` and `object/boolean.go` each define their own
private `isTruthy(Object) bool` with identical logic (`object/boolean.go`'s
copy is exported as `IsTrue`/`IsFalse`). Not a bug — both copies agree
today — but two independent definitions of a rule this central (§8.5) is a
maintenance hazard the moment one of them is edited and the other is not.

---

## 14. Decisions for 1.0

The following are this specification's working answers to the
product-direction calls §12 and §13 raise but can't resolve on their own —
made here so 1.0 ships with an answer rather than an asterisk. Revisit any of
them only if real usage argues otherwise; until then, this is what 1.0
targets.

1. **User-defined function arity.** Functions and methods defined in Ghost
   get the same strict arity checking every library function already has
   (§8.7, §12). A call with the wrong number of arguments becomes an
   `Argument` fault naming the call, exactly as it already does for a
   library method — closing the largest remaining behavioral gap between
   user code and library code, and the one most in tension with the "no
   silent gaps" goal in §3.
2. **`Map` iteration order.** `Map` guarantees insertion order for
   `keys()`, `values()`, `for ... in`, and `String()` (§13.5), matching the
   predictability JS objects and PHP associative arrays already give their
   users (§2). This likely requires backing `Map` with an insertion-ordered
   structure rather than a bare Go map — a larger implementation change, but
   the goal is fixed regardless of how large the fix turns out to be.
3. **`string.find`/`findAll`/`matches`.** These flip to the conventional
   `subject.find(pattern)` shape before 1.0 locks the API (§13.7), matching
   JS, PHP, and Python, and the "expressive, familiar syntax" goal in §1.
   `findAll` is fixed at the same time to return every match, not just the
   first.
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
