package evaluator

import (
	"fmt"
	"sort"
	"strings"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// undefined reports a name that is not bound anywhere the code could see it.
//
// Almost every one of these is a typo, and the interpreter is holding the list
// of names that were in scope at the moment it gave up. Offering the nearest of
// them turns the most common error in the language from a statement of fact
// into something a reader can act on without going looking.
func undefined(tok token.Token, name string, scope *object.Scope) *object.Error {
	raised := object.NewError(fault.Name, tok, "`%s` is not defined", name)

	// A name that is exactly a registered module or function — the standard
	// library's, or one an embedder registered under a scheme of its own
	// (library.RegisterModuleForScheme) — is probably not a typo at all. That
	// is a different mistake from a misspelling, and deserves a different
	// help line: the import to write, not a "did you mean" pointing at
	// itself.
	if schemes := library.Locate(name); len(schemes) > 0 {
		raised.WithHelp("`%s` is not a global — import it: %s", name, importHints(schemes, name))

		return raised
	}

	// A name that is a member of the class this code is running in is not a
	// missing variable, it is a missing receiver. Methods are reached through
	// `this` (§8.8, §14 decision 12) and deliberately do not sit in the
	// lexical chain, so a bare sibling call lands here - and this is the one
	// place that can say why (§13.17).
	if scope != nil && scope.Class != nil {
		if _, _, ok := object.LookupMember(scope.Class, name); ok {
			raised.WithHelp("`%s` is a member of `%s`; reach it through `this.%s`", name, scope.Class.Name.Value, name)

			return raised
		}
	}

	if suggestion, ok := nearestName(name, visibleNames(scope)); ok {
		raised.WithHelp("did you mean `%s`?", suggestion)

		return raised
	}

	// Nothing in scope was close, so try a near-miss on a registered name
	// too — `mathh.pi` should still point at `math`, just with the import it
	// now needs rather than a bare "did you mean".
	if suggestion, ok := nearestName(name, stdlibNames()); ok {
		raised.WithHelp("did you mean `%s`? import it: %s", suggestion, importHints(library.Locate(suggestion), suggestion))
	}

	return raised
}

// importHints renders the `import "scheme:name"` to add, or — on the rare
// occasion more than one scheme registers the same name — every scheme it
// would have to be disambiguated between.
func importHints(schemes []string, name string) string {
	hints := make([]string, len(schemes))

	for index, scheme := range schemes {
		hints[index] = fmt.Sprintf("`import \"%s:%s\"`", scheme, name)
	}

	return strings.Join(hints, " or ")
}

// stdlibNames lists every module and function registered under any scheme —
// the standard library's own and every embedder's — for suggesting a
// near-miss once no in-scope name is close enough.
func stdlibNames() []string {
	names := make([]string, 0, 32)

	for _, scheme := range library.Schemes() {
		registry := library.Scheme(scheme)

		for name := range registry.Modules {
			names = append(names, name)
		}

		for name := range registry.Functions {
			names = append(names, name)
		}
	}

	return names
}

// tooDeep reports a recursion that never bottomed out. It is the one error in
// Ghost that is about the interpreter's own limits rather than the program's
// values, and it says so plainly, because the fix is always in the program.
func tooDeep(tok token.Token) *object.Error {
	return object.NewError(fault.Value, tok, "call depth of %d exceeded", maxCallDepth).
		WithHelp("a function is almost certainly calling itself without ever stopping; check its base case")
}

// unknownMethod reports a method a class does not have, suggesting the nearest
// one it does. A misremembered method name is the same mistake as a misspelled
// variable, and deserves the same help.
func unknownMethod(tok token.Token, class *object.Class, name string) *object.Error {
	raised := object.NewError(fault.Property, tok, "class `%s` has no method `%s`", class.Name.Value, name)

	if suggestion, ok := nearestName(name, classMembers(class)); ok {
		raised.WithHelp("did you mean `%s`?", suggestion)
	}

	return raised
}

// classMembers lists every name a class responds to, its ancestors and the
// traits they use included.
func classMembers(class *object.Class) []string {
	names := make([]string, 0, 16)

	for _, ancestor := range class.Ancestors() {
		for name := range ancestor.Environment.All() {
			names = append(names, name)
		}

		for _, field := range ancestor.Fields {
			names = append(names, field.Name)
		}

		for _, trait := range ancestor.Traits {
			for name := range trait.Environment.All() {
				names = append(names, name)
			}
		}
	}

	return names
}

// moduleSuggestion offers the nearest method a module actually has. It answers
// with an empty string when nothing is close, and an empty help line is left
// out of the report.
func moduleSuggestion(module *object.LibraryModule, name string) string {
	names := make([]string, 0, len(module.Methods))

	for method := range module.Methods {
		names = append(names, method)
	}

	suggestion, ok := nearestName(name, names)

	if !ok {
		return ""
	}

	return "did you mean `" + module.Name + "." + suggestion + "()`?"
}

// modulePropertySuggestion offers the nearest property a module actually has.
func modulePropertySuggestion(module *object.LibraryModule, name string) string {
	names := make([]string, 0, len(module.Properties)+len(module.Classes))

	for property := range module.Properties {
		names = append(names, property)
	}

	for class := range module.Classes {
		names = append(names, class)
	}

	suggestion, ok := nearestName(name, names)

	if !ok {
		return ""
	}

	return "did you mean `" + module.Name + "." + suggestion + "`?"
}

// operatorError reports an operator applied to operands it has no meaning for.
//
// Two shapes of that mistake read very differently — mixing a number with a
// string is not the same as asking for an ordering between two lists — so the
// message names the operands when they differ and names the type once when
// they do not.
func operatorError(tok token.Token, operator token.Type, left object.Object, right object.Object) *object.Error {
	if left.Type() == right.Type() {
		return object.NewError(fault.Type, tok, "cannot use `%s` between two %s", operator, plural(object.TypeName(left)))
	}

	return object.NewError(fault.Type, tok, "cannot use `%s` between %s and %s", operator, object.TypeName(left), object.TypeName(right))
}

// logicalOperandError reports a non-boolean operand of `and`/`or`. It names
// the side at fault rather than both types, because short-circuiting means
// there is not always another side to name: when the left operand is wrong
// the right one has deliberately not been evaluated, and reporting a type it
// might have had would be inventing one.
//
// The null case gets the help line because it is the mistake this wording
// exists for - a guard written as `x and x.field`, in the idiom of a language
// where `and` is truthy rather than boolean.
func logicalOperandError(tok token.Token, operator token.Type, operand object.Object, side string) *object.Error {
	raised := object.NewError(fault.Type, tok, "cannot use `%s` with %s on the %s", operator, object.TypeName(operand), side)

	if object.TypeName(operand) == "null" {
		raised.WithHelp("compare it first, as in `x != null`")
	}

	return raised
}

// memberCollisionError rejects one name declared as both a field and a method
// in the same class or trait body (§13.18). Both declarations are legal on
// their own and neither lookup path knows the other exists - a property read
// answers with the field, a call answers with the method - so the only place
// the collision can be reported is where the second declaration creates it.
func memberCollisionError(tok token.Token, name string, existing string) *object.Error {
	return object.NewError(fault.Syntax, tok, "`%s` is already declared as a %s in this body", name, existing).
		WithHelp("rename one of them; otherwise `%s` reads the field and `%s()` calls the method, which is not a difference a reader will predict", name, name)
}

// plural writes a type name as it reads when there is more than one of them.
func plural(name string) string {
	if strings.HasSuffix(name, "s") || strings.HasSuffix(name, "x") || strings.HasSuffix(name, "ch") {
		return name + "es"
	}

	return name + "s"
}

// indexTypeError reports a subscript applied to something that cannot take one.
// The two operands fail for different reasons — a list indexed by a string is a
// different mistake from a boolean indexed at all — so the message says which
// of them is the problem.
func indexTypeError(node *ast.Index, left object.Object, index object.Object) *object.Error {
	switch left.Type() {
	case object.LIST, object.STRING:
		return object.NewError(fault.Type, node.Token, "a %s index has to be a number, got %s", object.TypeName(left), object.TypeName(index))
	}

	return object.NewError(fault.Type, node.Token, "cannot index %s", object.TypeName(left)).
		WithHelp("indexing works on lists, maps, and strings")
}

// visibleNames collects every name the code at this point could have meant: the
// bindings in scope, walking outwards, and the names reachable without an
// import (console, type). The rest of the standard library is deliberately
// left out here — a near-miss on an import-only name is handled by undefined()
// itself, which offers the import to write rather than pretending the bare
// name would resolve.
func visibleNames(scope *object.Scope) []string {
	names := make([]string, 0, 32)

	if scope != nil {
		for name := range scope.Environment.All() {
			names = append(names, name)
		}
	}

	for name := range library.Functions {
		if library.IsGlobal(name) {
			names = append(names, name)
		}
	}

	for name := range library.Modules {
		if library.IsGlobal(name) {
			names = append(names, name)
		}
	}

	return names
}

// nearestName picks the candidate closest to a misspelling, if one is close
// enough to be worth mentioning. The threshold scales with the length of the
// name: a two-letter name has to match almost exactly, while a longer one can
// survive a couple of slips. Guessing wrongly is worse than not guessing, so
// the bar is deliberately high.
func nearestName(name string, candidates []string) (string, bool) {
	limit := len(name) / 3

	if limit < 1 {
		limit = 1
	}

	if limit > 3 {
		limit = 3
	}

	best := ""
	bestDistance := limit + 1

	// Sorting keeps the answer stable: map iteration is random, and two
	// candidates the same distance away would otherwise alternate between runs.
	sort.Strings(candidates)

	for _, candidate := range candidates {
		if candidate == name {
			continue
		}

		distance := editDistance(strings.ToLower(name), strings.ToLower(candidate))

		if distance > bestDistance {
			continue
		}

		if distance < bestDistance || sharesMore(name, candidate, best) {
			best = candidate
			bestDistance = distance
		}
	}

	return best, best != ""
}

// sharesMore breaks a tie between two candidates the same distance from a
// misspelling. The one that agrees with it for longer wins, so `pii` is
// answered with `pi` rather than with `phi`, which is equally close by the
// arithmetic and obviously further away to a reader.
func sharesMore(name string, candidate string, best string) bool {
	if best == "" {
		return true
	}

	return commonPrefix(name, candidate) > commonPrefix(name, best)
}

// commonPrefix counts the characters two names begin with in common.
func commonPrefix(left string, right string) int {
	shared := 0

	for shared < len(left) && shared < len(right) && left[shared] == right[shared] {
		shared++
	}

	return shared
}

// editDistance is the Damerau-Levenshtein distance between two strings: the
// number of single-character insertions, deletions, substitutions, or
// transpositions that turn one into the other.
//
// Transpositions are counted because they are what typing actually produces.
// Plain Levenshtein scores `nmae` against `name` as two changes and would put
// it out of reach of any threshold tight enough to be trusted; counting the
// swap as the one slip it was brings the suggestion back.
func editDistance(from string, to string) int {
	source := []rune(from)
	target := []rune(to)

	if len(source) == 0 {
		return len(target)
	}

	if len(target) == 0 {
		return len(source)
	}

	// A transposition needs the row before last, so three rows are kept rather
	// than the whole grid.
	beforeLast := make([]int, len(target)+1)
	previous := make([]int, len(target)+1)
	current := make([]int, len(target)+1)

	for column := range previous {
		previous[column] = column
	}

	for row := 1; row <= len(source); row++ {
		current[0] = row

		for column := 1; column <= len(target); column++ {
			cost := 1

			if source[row-1] == target[column-1] {
				cost = 0
			}

			current[column] = smallest(
				current[column-1]+1,
				previous[column]+1,
				previous[column-1]+cost,
			)

			if row > 1 && column > 1 && source[row-1] == target[column-2] && source[row-2] == target[column-1] {
				current[column] = smallest(current[column], beforeLast[column-2]+1)
			}
		}

		beforeLast, previous, current = previous, current, beforeLast
	}

	return previous[len(target)]
}

func smallest(values ...int) int {
	least := values[0]

	for _, value := range values[1:] {
		if value < least {
			least = value
		}
	}

	return least
}
