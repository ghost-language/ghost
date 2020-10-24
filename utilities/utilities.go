package utilities

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ghostlang.org/x/ghost/object"
)

var SearchPaths []string

func init() {
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("error getting cwd: %s", err)
	}

	if e := os.Getenv("GHOSTPATH"); e != "" {
		tokens := strings.Split(e, ":")

		for _, token := range tokens {
			AddPath(token)
		}
	} else {
		SearchPaths = append(SearchPaths, cwd)
	}
}

func AddPath(path string) error {
	path = os.ExpandEnv(filepath.Clean(path))
	absolutePath, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	SearchPaths = append(SearchPaths, absolutePath)

	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func FindModule(name string) string {
	basename := fmt.Sprintf("%s.ghost", name)

	for _, p := range SearchPaths {
		filename := filepath.Join(p, basename)

		if Exists(filename) {
			return filename
		}
	}

	return ""
}

// NewError returns a new error object used during runtime.
func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
