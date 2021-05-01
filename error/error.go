package error

import (
	"fmt"

	"ghostlang.org/x/ghost/object"
)

// Define all error constants later used to define error messages.
const (
	_ int = iota
	ArgumentMustBe
	ArgumentNotSupported
	ErrorReadingInput
	ErrorReadingModule
	IndexOutOfRange
	InvalidImportPath
	InvalidPropertyType
	InvalidPropertyAssignment
	InfixTypeMismatch
	NoModuleFound
	NotAFunction
	ParseError
	UnknownOperator
	UnknownIdentifier
	UnknownInfixOperator
	UnsupportedIndexExpression
	UnsupportedPostfixOperator
	UnsupportedPrefixOperator
	UnusableForLoop
	UnusableMapKey
	WrongNumberArguments
)

var MessageBag = map[int]string{
	ArgumentMustBe:             "%s argument to '%s' must be %s, got %s",
	ArgumentNotSupported:       "%s argument to '%s' is not supported, got %s",
	ErrorReadingInput:          "error reading input: %s",
	ErrorReadingModule:         "error reading module '%s': %s",
	IndexOutOfRange:            "index out of range: %d",
	InvalidImportPath:          "invalid import path: %s",
	InvalidPropertyType:        "invalid property '%s' on type %s",
	InvalidPropertyAssignment:  "invalid property assignment '%s' on type %s",
	InfixTypeMismatch:          "type mismatch: %s %s %s",
	NoModuleFound:              "no module named '%s' found",
	NotAFunction:               "not a function: %s",
	ParseError:                 "parse error: %s",
	UnknownOperator:            "unknown operator: %s%s",
	UnknownIdentifier:          "unknown identifier: %s",
	UnknownInfixOperator:       "unknown operator: %s %s %s",
	UnsupportedIndexExpression: "index expression not supported on type %s",
	UnsupportedPostfixOperator: "unsupported operator for postfix expression: '%s' and type: %s",
	UnsupportedPrefixOperator:  "unsupported operator for prefix expression: '%s' and type: %s",
	UnusableForLoop:            "unusable as for loop: %s",
	UnusableMapKey:             "unusable as map key: %s",
	WrongNumberArguments:       "wrong number of arguments: %d while expected: %d",
}

// NewError formats and returns a new Error object.
func NewError(line int, index int, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(MessageBag[index], args...) + fmt.Sprintf(" on line %d", line)}
}

// IsError determines if the passed object is an error object.
func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
