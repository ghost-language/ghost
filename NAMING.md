# Naming

This is the convention every name in Ghost's language and standard library
follows: keywords, library functions and modules, and the methods defined on
each object type. It was written down after an audit found the convention
already in consistent use but nowhere stated - this is that statement, so a
new addition has a rule to follow instead of just precedent to imitate.

Ghost's syntax already commits to JavaScript/TypeScript conventions for
classes (`class`, `extends`, `new`), so naming follows that lead rather than
reaching for PHP's `snake_case`. What's borrowed from Laravel here is the
discipline of a written convention and comprehensive, predictable coverage of
the standard operations a language's users reach for - not PHP's casing.

## The rules

**Names are `camelCase`.** Functions, methods, module names, and properties -
`toLowerCase`, `randomInt`, `isPrime`. Classes are `StudlyCase`, matching the
JS/TS convention the class syntax already follows.

**Actions are verbs.** `push`, `pop`, `split`, `trim`, `reshape` - not
`pusher`, `splitting`, `trimmed`. A method name is the imperative form of what
it does.

**Predicates ask a yes/no question and read like one.** Every boolean-
returning library function is named `isX` or `hasX` - `isNaN`, `isEven`,
`isPrime`, `isFinite`, `hasA`. If a name doesn't read as a question, it
shouldn't return a boolean.

**Conversions read left to right, target last.** `toString`, `toNumber`,
`toLowerCase`, `toUpperCase` - the type or form being converted *to* is always
the suffix, so the name reads the same direction as the conversion happens.
`toString()` in particular is universal: every value type answers it,
booleans and null included, so a reader never has to check whether a
particular type happens to support turning itself into a string.

**Modules are lowercase domains, not types.** `math`, `date`, `random`,
`console`, `os` name a *place* an operation belongs, not a kind of value.
Before adding a method to a module, ask which domain the operation is really
in - randomness belongs in `random`, not `math`; a calendar instant belongs in
`date`, not `os`; pausing or ending the program itself belongs in `os`, not
`date`, which is why `os.sleep()` lives beside `os.exit()` rather than beside
`date.now()` - rather than adding a method to whichever module is already
open.

**A name means one thing.** If two modules define a method with the same
name, they'd better mean the exact same operation with the exact same
arguments - the way `a + b` and `math.add(a, b)` are one operation reached two
ways (see `object.Broadcast`). If the arguments or the answer differ, the
names have to differ too, however closely related the operations are.
`math.randomInt()` and `random.random()` are the example this rule is named
for: both draw from the same generator, but one always answers a whole number
and the other a float, so they are not the same name.

**A getter and its setter don't share a name.** A property that can also be
set through a method call needs two different words - one to read it, one to
change it (`random.currentSeed` the property, `random.seed()` the method) -
so that which one a piece of code means is never a guess from context alone.

**Errors describe, they don't locate.** This is a naming discipline as much
as an error-message one: `` "`%s` is not defined" ``, not a sentence built
around a position. See the "Error handling" section of CLAUDE.md for the full
rules - they're the same spirit applied to sentences instead of identifiers.

## What counts as a genuine collision

Not every case where two things could theoretically share a word is a
problem. A rename or a consolidation is worth making when:

- Two names mean different things a reader can't tell apart without opening
  both definitions (this is what happened to `math.random`/`random.random`
  before the audit - same word, incompatible argument shapes).
- A name's contract is silently different from what an identical name means
  elsewhere in the language (console's output methods bypassing the writer
  `print()` respects, before the audit fixed it).
- A capability lives in the wrong domain entirely (current-time precision
  living partly in `os`, before it moved to `date`).

It is not a problem when two access points deliberately share an
implementation and agree completely on what they mean - that's the same
pattern as arithmetic operators and the math module's methods, and it doesn't
need fixing, only documenting.

## Breaking changes

Ghost has no deprecation mechanism today: a renamed method simply stops
existing under its old name in the next version. A rename is therefore a
breaking change for any `.ghost` script written against the old name, made
deliberately rather than as a side effect of a naming pass, and called out in
the version's release notes.
