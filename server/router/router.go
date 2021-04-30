package router

import (
	"net/http"
	"regexp"
)

var Routes = map[string]*Route{}

type Route struct {
	Method string
	Regex *regexp.Regexp
	Handler http.HandlerFunc
}

// RegisterRoute ...
func RegisterRoute(method, pattern string, handler http.HandlerFunc) {
	Routes[pattern] = &Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}