package version

import "fmt"

var (
	Version = "v0.1.1"
)

func String() string {
	return fmt.Sprintf("%s", Version)
}
