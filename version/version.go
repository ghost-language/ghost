package version

import "fmt"

// Version signifies the current version of Ghost.
const Version = "experimental"

// String returns the version formatted as a string.
func String() string {
	return fmt.Sprintf("%s", Version)
}
