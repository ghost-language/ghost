package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/lexer"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/version"
)

// PROMPT designates the REPL prompt characters to accept
// user input.
const PROMPT = ">> "

// OUTPUT designates the REPL output characters to display
// program results.
const OUTPUT = "   "

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Options struct {
	Interactive bool
	Server bool
}

type REPL struct {
	args []string
	opts *Options
}

func New(args []string, opts *Options) *REPL {
	return &REPL{args, opts}
}

func (r *REPL) Run() {
	registerCloseHandler()

	if len(r.args) == 0 {
		fmt.Printf("Ghost (%s)\n", version.String())
		fmt.Printf("Press Ctrl + C to exit\n\n")

		r.StartEvalLoop(os.Stdin, os.Stdout, nil)
		return
	}

	if len(r.args) > 0 {
		if r.opts.Server {
			r.StartServer()
		} else {
			start := time.Now()
			f, err := os.Open(r.args[0])

			if err != nil {
				log.Fatalf("Could not open source file %s: %s", r.args[0], err)
			}

			defer f.Close()

			env := r.Eval(f, os.Stdout)
			elapsed := time.Since(start)

			if r.opts.Interactive {
				fmt.Printf(OUTPUT+"(executed in: %s)\n", elapsed)
				r.StartEvalLoop(os.Stdin, os.Stdout, env)
			}
		}
	}
}

func (r *REPL) Eval(f io.Reader, writer io.Writer) (env *object.Environment) {
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
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	obj := evaluator.Eval(program, env)

	if obj != nil {
		if _, ok := obj.(*object.Error); ok {
			io.WriteString(os.Stdout, OUTPUT+obj.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	}

	return
}

func (r *REPL) StartEvalLoop(in io.Reader, out io.Writer, env *object.Environment) {
	scanner := bufio.NewScanner(in)

	if env == nil {
		env = object.NewEnvironment()
	}

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		obj := evaluator.Eval(program, env)

		if obj != nil {
			if _, ok := obj.(*object.Null); !ok {
				io.WriteString(out, OUTPUT+obj.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

func (r *REPL) StartServer() {
	address := "0.0.0.0:8080"
	currentTime := time.Now()

	fmt.Printf("%s --> ", currentTime.Format("2006/01/02 15:04:05"))
	fmt.Printf(InfoColor, fmt.Sprintf("Starting Ghost %s server: ", version.Version))
	fmt.Printf("%s\n", address)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		f, err := os.Open(r.args[0])

		if err != nil {
			log.Fatalf("Could not open source file %s: %s", r.args[0], err)
		}

		defer f.Close()

		r.Eval(f, writer)

		secs := time.Since(start).String()

		log.Printf("--> %s %s (%s)", request.Method, request.URL.Path, secs)
	})

	log.Fatal(http.ListenAndServe(address, nil))
}

func registerCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Printf("\nExiting...\n")
		os.Exit(0)
	}()
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "\tPARSE ERROR:\n")

	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
