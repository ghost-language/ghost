package evaluator

import "testing"

// The tests in this file cover §13.18 and §13.22 — the two ways a class or
// trait member can collide with a name that already means something, both
// reported where the collision is created rather than left to surface as a
// confusing failure elsewhere.
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

func TestMethodShadowingAnImportedModuleIsRejected(t *testing.T) {
	t.Run("a scheme import shadowed by a method of the same name", func(t *testing.T) {
		input := "import \"ghost:math\" as math\nclass Theme {\n\tmath(role) { return role }\n}"

		isErrorObject(t, evaluate(input), "test.gs:3:2: syntax error: method `math` shadows an imported module of the same name")
	})

	t.Run("the same shadowing inside a trait", func(t *testing.T) {
		input := "import \"ghost:math\" as math\ntrait Sums {\n\tmath() { return 1 }\n}"

		isErrorObject(t, evaluate(input), "test.gs:3:2: syntax error: method `math` shadows an imported module of the same name")
	})

	t.Run("a Ghost module bound as a map is shadowed just the same", func(t *testing.T) {
		dir := t.TempDir()
		writeModule(t, dir, "geometry.gs", `origin = 0`)

		result := evaluateInDirectory(dir, "import \"geometry\" as geometry\nclass A {\n\tgeometry() { return 1 }\n}")

		isErrorObject(t, result, "main.gs:3:2: syntax error: method `geometry` shadows an imported module of the same name")
	})
}

// TestMethodShadowingAllowsWhatIsNotAModule is the other half of §13.22's
// check: it fires on imported modules and nothing else. A method is allowed to
// carry the same name as any ordinary outer binding, which is ordinary lexical
// shadowing and not the hazard — and the aliasing workaround Chisel adopted
// has to keep working, since it is what the fix tells people to do.
func TestMethodShadowingAllowsWhatIsNotAModule(t *testing.T) {
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
			name:     "a method named like a plain map, which is not a module",
			input:    "config = {a: 1}\nclass A { config() { return 2 } }\nnew A().config()",
			expected: 2,
		},
		{
			name:     "importing under a different name keeps both reachable",
			input:    "import \"ghost:math\" as mathModule\nclass A {\n\tmath() { return 1 }\n\tfloor() { return mathModule.floor(2.7) }\n}\nnew A().floor()",
			expected: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isNumberObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestGlobalModulesAreNotShadowedByMethods pins the asymmetry that makes the
// §13.22 check correct to scope to imports. A global module is resolved
// through the library registry rather than the environment chain, so a method
// of the same name does not hide it and there is nothing to report.
func TestGlobalModulesAreNotShadowedByMethods(t *testing.T) {
	input := "class A {\n\tconsole() { return 1 }\n\tprobe() { return console.log }\n}\nnew A().console()"

	isNumberObject(t, evaluate(input), 1)
}
