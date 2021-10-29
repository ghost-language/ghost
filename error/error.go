package error

// Reasons are enumerated here to be used in the Error struct. Defining
// reasons here keeps things streamlined and consistent through the codebase.
const (
	Unknown = 0
	Syntax  = 1
	Runtime = 2
	System  = 50
)

// Error struct represents errors encountered during the scanning, parsing, and
// evaluation phases of Ghost.
type Error struct {
	Reason  int
	Message string
}
