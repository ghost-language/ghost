package object

import "fmt"

const LIBRARY_PROPERTY = "LIBRARY_PROPERTY"

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
func (libraryProperty *LibraryProperty) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
