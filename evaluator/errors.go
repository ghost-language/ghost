package evaluator

import (
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

	if suggestion, ok := nearestName(name, visibleNames(scope)); ok {
		raised.WithHelp("did you mean `%s`?", suggestion)
	}

	return raised
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
	names := make([]string, 0, len(module.Properties))

	for property := range module.Properties {
		names = append(names, property)
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
// bindings in scope, walking outwards, and the library globals.
func visibleNames(scope *object.Scope) []string {
	names := make([]string, 0, 32)

	if scope != nil {
		for name := range scope.Environment.All() {
			names = append(names, name)
		}
	}

	for name := range library.Functions {
		names = append(names, name)
	}

	for name := range library.Modules {
		names = append(names, name)
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
