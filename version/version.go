package version

import "fmt"

var (
	Version = "dev-nightly"
)

func String() string {
	return fmt.Sprintf("%s", Version)
}
