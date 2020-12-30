package error

import (
	"fmt"

	"ghostlang.org/x/ghost/object"
)

// Define all error constants later used to define error messages.
const (
	_ int = iota
	Placeholder
	ArgumentMustBe
	ErrorReadingModule
	InfixTypeMismatch
	NoModuleFound
	ParseError
	UnknownOperator
	UnknownIdentifier
	UnknownInfixOperator
	UnsupportedPostfixOperator
	UnsupportedPrefixOperator
	WrongNumberArguments
)

var messageBag = map[int]string{
	Placeholder:                "placeholder error message",
	ArgumentMustBe:             "%s argument must be %s",
	ErrorReadingModule:         "error reading module '%s': %s",
	InfixTypeMismatch:          "type mismatch: %s %s %s",
	NoModuleFound:              "no module named '%s' found",
	ParseError:                 "parse error: %s",
	UnknownOperator:            "unknown operator: %s%s",
	UnknownIdentifier:          "unknown identifier: %s",
	UnknownInfixOperator:       "unknown operator: %s %s %s",
	UnsupportedPostfixOperator: "unsupported operator for postfix expression: '%s' and type: %s",
	UnsupportedPrefixOperator:  "unsupported operator for prefix expression: '%s' and type: %s",
	WrongNumberArguments:       "wrong number of arguments: got=%d, expected=%d",
}

// NewError formats and returns a new Error object.
func NewError(line int, index int, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(messageBag[index], args...) + fmt.Sprintf(" on line %d", line)}
}

// IsError determines if the passed object is an error object.
func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
