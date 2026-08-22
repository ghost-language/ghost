package object

import (
	"fmt"

	"ghostlang.org/x/ghost/token"
)

// LibraryProperty objects consist of a native Go property.
type LibraryProperty struct {
	Name     string
	Property GoProperty
}

// String represents the library property's value as a string.
func (libraryProperty *LibraryProperty) String() string {
	return fmt.Sprintf("library property {%s}", libraryProperty.Name)
}

// Type returns the library property object type.
func (libraryProperty *LibraryProperty) Type() Type {
	return LIBRARY_PROPERTY
}

// Method defines the set of methods available on library property objects.
func (libraryProperty *LibraryProperty) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}
