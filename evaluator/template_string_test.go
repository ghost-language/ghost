package evaluator

import "testing"

// Template literals stringify each interpolated value the same way every
// other native stringification point does (console.log, print,
// string.format): by calling the object's own String() representation, so no
// type needs an explicit toString() call to be interpolated.
func TestTemplateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no interpolation", "`hello`", "hello"},
		{"empty template", "``", ""},
		{"a number needs no toString", "`count: ${41 + 1}`", "count: 42"},
		{"a boolean needs no toString", "`ready: ${true}`", "ready: true"},
		{"a list needs no toString", "`items: ${[1, 2, 3]}`", "items: [1, 2, 3]"},
		{"null interpolates as the word null", "`value: ${null}`", "value: null"},
		{"multiple interpolations", "name = \"Fido\" age = 3 `${name} is ${age}`", "Fido is 3"},
		{"an interpolation can itself be a map literal", "`${ {\"a\": 1}[\"a\"] }`", "1"},
		{"a nested template literal", "`outer ${ `inner ${1 + 1}` }`", "outer inner 2"},
		{"an escaped dollar brace is not an interpolation", "`price: \\${5}`", "price: ${5}"},
		{"a backtick can be escaped inside a template", "`it\\`s here`", "it`s here"},
		{"a template can span multiple lines", "`a\nb`", "a\nb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)

			isStringObject(t, result, tt.expected)
		})
	}
}

// A runtime error raised while evaluating an interpolated expression surfaces
// exactly as it would outside a template, positioned at the expression that
// actually failed.
func TestTemplateStringPropagatesErrors(t *testing.T) {
	result := evaluate("`bad: ${1 + true}`")

	isErrorObject(t, result, "test.gs:1:11: type error: cannot use `+` between number and boolean")
}

// Interpolating a class instance falls back to the same description
// console.log already gives one, since neither goes through a user-defined
// toString() method.
func TestTemplateStringInterpolatesInstancesLikeConsoleLog(t *testing.T) {
	result := evaluate("class Dog {} d = new Dog() `dog: ${d}`")

	isStringObject(t, result, "dog: class instance Dog")
}
