package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/object"
)

// The tests in this file cover §13.18, §13.22 and §13.17 — what happens when a
// class or trait member shares a name with something else.
//
// §13.18, a field and a method sharing one name, is a genuine collision and is
// rejected where it is written: the two used to coexist in silence, a property
// read answering with the field and a call answering with the method.
//
// §13.22 and §13.17 are the same question answered once by §14 decision 12: a
// method is a member reached through `this`, not a lexical binding. So a method
// never enters the scope chain, cannot shadow an import (§13.22), and a bare
// sibling call is a plain name error rather than a call with the wrong
// receiver (§13.17).
//
// §13.18 is a field and a method sharing one name, which used to coexist in
// silence: `x.thing` answered with the field and `x.thing()` with the method,
// because the two lookup paths never meet.
//
// §13.22 is a method shadowing an imported module, which used to surface as a
// property error at whatever call site first wanted the module, with nothing
// at the import or the method to connect it back.

func TestMemberCollisionIsRejected(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "a method colliding with a field declared before it",
			input:    "class Thing {\n\tthing = 7\n\tthing() { return 1 }\n}",
			expected: "test.gs:3:2: syntax error: `thing` is already declared as a field in this body",
		},
		{
			name:     "a field colliding with a method declared before it",
			input:    "class Thing {\n\tthing() { return 1 }\n\tthing = 7\n}",
			expected: "test.gs:3:2: syntax error: `thing` is already declared as a method in this body",
		},
		{
			name:     "the same collision in a trait body",
			input:    "trait Loud {\n\tvolume = 3\n\tvolume() { return 4 }\n}",
			expected: "test.gs:3:2: syntax error: `volume` is already declared as a field in this body",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isErrorObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestMemberCollisionAllowsWhatIsNotOne guards the check against the shapes it
// must not reject. Each of these was legal before and stays legal.
func TestMemberCollisionAllowsWhatIsNotOne(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "a subclass overriding an inherited method",
			input:    "class A { speak() { return 1 } }\nclass B extends A { speak() { return 2 } }\nnew B().speak()",
			expected: 2,
		},
		{
			name:     "a subclass field overriding an inherited field",
			input:    "class A { size = 1 }\nclass B extends A { size = 2 }\nnew B().size",
			expected: 2,
		},
		{
			name:     "a field and a method that merely sit beside each other",
			input:    "class A { size = 1\n\tdouble() { return this.size * 2 } }\nnew A().double()",
			expected: 2,
		},
		{
			name:     "a method named like a field of an unrelated class",
			input:    "class A { thing = 1 }\nclass B { thing() { return 2 } }\nnew B().thing()",
			expected: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isNumberObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestAMethodDoesNotShadowAnImport is §13.22, resolved by scoping rather than
// by a diagnostic (§14 decision 12). A method is a member, so it never enters
// the lexical chain and cannot hide a name from its siblings: inside the class
// `math` is the import and `this.math()` is the method, and both work.
func TestAMethodDoesNotShadowAnImport(t *testing.T) {
	t.Run("a method may share a name with a scheme import", func(t *testing.T) {
		input := `
			import "ghost:math" as math
			class Theme {
				load() { return math.floor(3.7) }
				math(role) { return role }
			}
			new Theme().load()
		`

		isNumberObject(t, evaluate(input), 3)
	})

	t.Run("and the method of that name is still reachable", func(t *testing.T) {
		input := `
			import "ghost:math" as math
			class Theme {
				load() { return math.floor(3.7) }
				math(role) { return role }
			}
			new Theme().math("accent")
		`

		isStringObject(t, evaluate(input), "accent")
	})

	t.Run("the same holds in a trait", func(t *testing.T) {
		input := `
			import "ghost:math" as math
			trait Sums {
				math() { return 1 }
				total() { return math.floor(2.5) }
			}
			class Adder { use Sums }
			new Adder().total()
		`

		isNumberObject(t, evaluate(input), 2)
	})

	t.Run("a field initializer resolves the same way a method body does", func(t *testing.T) {
		// This path is separate from the method one and would otherwise keep
		// the old behavior: the initializer is evaluated per instance, and
		// used to run in the class's member table.
		input := `
			import "ghost:math" as math
			class A {
				size = math.floor(2.7)
				math() { return 1 }
			}
			new A().size
		`

		isNumberObject(t, evaluate(input), 2)
	})

	t.Run("a Ghost source module is not shadowed either", func(t *testing.T) {
		dir := t.TempDir()
		writeModule(t, dir, "geometry.gs", `origin = 7`)

		result := evaluateInDirectory(dir, `
			import "geometry" as geometry
			class A {
				geometry() { return 1 }
				read() { return geometry.origin }
			}
			new A().read()
		`)

		isNumberObject(t, result, 7)
	})
}

// TestBareSiblingCallsAreNotDefined covers §13.17. §8.8 used to promise that a
// sibling method could be called by bare name; it resolved but ran with the
// wrong receiver, failing only if the callee touched `this` — so the failure
// depended on the callee's body and was reported at the callee rather than at
// the call. Methods now stay out of the lexical chain entirely, which makes
// the bare form an honest name error at the call site, with the receiver it is
// missing named in the help.
func TestBareSiblingCallsAreNotDefined(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "a sibling method that touches this",
			input:    "class W {\n\tconstructor(n) { this.n = n }\n\tdescribe() { return this.n }\n\tshow() { return describe() }\n}\nnew W(5).show()",
			expected: "test.gs:4:18: name error: `describe` is not defined",
		},
		{
			// This one used to "work", which is exactly why the old behavior
			// was so hard to reason about: whether a bare call succeeded
			// depended on whether the callee happened to use `this`.
			name:     "a sibling method that does not touch this",
			input:    "class W {\n\thelper() { return 42 }\n\tshow() { return helper() }\n}\nnew W().show()",
			expected: "test.gs:3:18: name error: `helper` is not defined",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isErrorObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestBareSiblingCallHelpNamesTheReceiver checks the help line, which is what
// turns this from a puzzle into a one-word fix.
func TestBareSiblingCallHelpNamesTheReceiver(t *testing.T) {
	result := evaluate("class W {\n\tdescribe() { return 1 }\n\tshow() { return describe() }\n}\nnew W().show()")

	err, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", result, result)
	}

	if err.Fault.Help != "`describe` is a member of `W`; reach it through `this.describe`" {
		t.Errorf("help has wrong text. got=%q", err.Fault.Help)
	}
}

// TestMethodsMayShareANameWithAnyOuterBinding is what is left of the old
// shadowing table: none of these were ever rejected, and none are now.
func TestMethodsMayShareANameWithAnyOuterBinding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "a method named like an outer function",
			input:    "function helper() { return 1 }\nclass A { helper() { return 2 } }\nnew A().helper()",
			expected: 2,
		},
		{
			name:     "a method named like an outer variable",
			input:    "label = 1\nclass A { label() { return 2 } }\nnew A().label()",
			expected: 2,
		},
		{
			name:     "a method named like a plain map",
			input:    "config = {a: 1}\nclass A { config() { return 2 } }\nnew A().config()",
			expected: 2,
		},
		{
			name:     "an outer variable stays readable from a method of the same name",
			input:    "label = 9\nclass A {\n\tlabel() { return 2 }\n\tread() { return label }\n}\nnew A().read()",
			expected: 9,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isNumberObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestAMethodLocalCannotClobberASiblingMethod pins a hazard §14 decision 9
// recorded as a live cost of walking assignment: a method's local temporary
// named like a sibling method used to rebind that method, since assignment
// walked into the class environment. Members are no longer in that chain
// (§14 decision 12), so the local stays a local.
func TestAMethodLocalCannotClobberASiblingMethod(t *testing.T) {
	input := `
		class A {
			scale() { return 1 }
			run() { scale = 5 return this.scale() }
		}
		new A().run()
	`

	isNumberObject(t, evaluate(input), 1)
}

// TestGlobalModulesAreNotShadowedByMethods keeps the global case pinned: a
// method named `console` never hid the global, and still does not.
func TestGlobalModulesAreNotShadowedByMethods(t *testing.T) {
	input := "class A {\n\tconsole() { return 1 }\n\tprobe() { return console.log }\n}\nnew A().console()"

	isNumberObject(t, evaluate(input), 1)
}
