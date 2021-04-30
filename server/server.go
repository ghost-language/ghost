package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/server/router"
	"ghostlang.org/x/ghost/version"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Server struct {
	args []string
}

func New(args []string) *Server {
	return &Server{args}
}

func (s *Server) Run() {
	address := "0.0.0.0:8080"
	currentTime := time.Now()

	fmt.Printf("%s --> ", currentTime.Format("2006/01/02 15:04:05"))
	fmt.Printf(InfoColor, fmt.Sprintf("Starting Ghost %s server: ", version.Version))
	fmt.Printf("%s\n", address)

	// http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	// 	status := "success"
	// 	start := time.Now()

	// 	f, err := os.Open(s.args[0])

	// 	if err != nil {
	// 		log.Fatalf("Could not open source file %s: %s", s.args[0], err)
	// 	}

	// 	defer f.Close()

	// 	_, errors := s.Eval(f, writer)

	// 	if len(errors) > 0 {
	// 		status = "error"

	// 		t, _ := template.ParseFiles("server/error.html")
	// 		// e := errors.ErrorBag{Message: errors.ParseErrorMessage}
	// 		t.Execute(writer, nil)
	// 	}

	// 	secs := time.Since(start).String()

	// 	log.Printf("--> %s (%s) %s (%s)", request.Method, status, request.URL.Path, secs)
	// })

	router := http.HandlerFunc(s.Router)

	log.Fatal(http.ListenAndServe(address, router))
}

func (s *Server) Eval(f io.Reader, writer io.Writer) (env *object.Environment, errors []string) {
	env = object.NewEnvironment()
	env.SetWriter(writer)

	b, err := ioutil.ReadAll(f)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source file: %s", err)
		return
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		// printParserErrors(os.Stderr, p.Errors())
		return
	}

	obj := evaluator.Eval(program, env)

	if obj != nil {
		if _, ok := obj.(*object.Error); ok {
			// io.WriteString(os.Stdout, OUTPUT+obj.Inspect())
			io.WriteString(os.Stdout, "error\n")
		}
	}

	return
}

func (s *Server) Router(writer http.ResponseWriter, request *http.Request) {
	status := "success"
	start := time.Now()
	var allow []string

	f, err := os.Open(s.args[0])

	if err != nil {
		log.Fatalf("Could not open source file %s: %s", s.args[0], err)
	}

	defer f.Close()

	_, errors := s.Eval(f, writer)

	if len(errors) > 0 {
		status = "error"

		t, _ := template.ParseFiles("server/error.html")
		// e := errors.ErrorBag{Message: errors.ParseErrorMessage}
		t.Execute(writer, nil)
	}

	fmt.Println("len:", len(router.Routes))

	for _, route := range router.Routes {
		matches := route.Regex.FindStringSubmatch(request.URL.Path)

		if len(matches) > 0 {
			if request.Method != route.Method {
				fmt.Println(request.Method + " != " + route.Method)
				allow = append(allow, route.Method)
				continue
			}

			ctx := context.WithValue(request.Context(), ctxKey{}, matches[1:])
			route.Handler(writer, request.WithContext(ctx))
			return
		}
	}

	if len(allow) > 0 {
		writer.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(writer, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(writer, request)

	secs := time.Since(start).String()

	log.Printf("--> %s (%s) %s (%s)", request.Method, status, request.URL.Path, secs)
}

type ctxKey struct{}