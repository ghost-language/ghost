package version

import "fmt"

var (
	Version = "0.5.0"
)

func String() string {
	return fmt.Sprintf("%s", Version)
}
