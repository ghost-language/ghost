package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"ghostlang.org/x/ghost/errors"
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/version"
	"github.com/peterh/liner"
	"github.com/spf13/viper"
)

// const ...
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

var (
	flagVersion     bool
	flagInteractive bool
	flagHelp        bool
	flagTokens      bool
	flagServe       bool
	flagAddress     string
	history         = filepath.Join(os.TempDir(), ".ghost_history")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&flagHelp, "h", false, "display help information")
	flag.BoolVar(&flagVersion, "v", false, "display version information")
	flag.BoolVar(&flagInteractive, "i", false, "enable interactive mode")
	flag.BoolVar(&flagTokens, "t", false, "output scanned token information")
	flag.BoolVar(&flagServe, "serve", false, "start web server")
	flag.StringVar(&flagAddress, "address", "0.0.0.0:8080", "listen on the given network address when running as a web server")
}

func main() {
	configureEnvironment()

	flag.Parse()

	if flagVersion {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), version.String())
		os.Exit(0)
	}

	if flagHelp {
		showHelp()
		os.Exit(2)
	}

	args := flag.Args()

	if len(args) > 1 {
		fmt.Println("Usage: ghost [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		if viper.GetString("SERVER") == "true" {
			runServer(args[0])
		} else {
			runFile(args[0])
		}
	} else {
		runPrompt()
	}
}

func configureEnvironment() {
	viper.SetConfigFile(".env")

	viper.ReadInConfig()

	viper.SetDefault("SERVER", false)
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
}

func runPrompt() {
	line := liner.NewLiner()
	env := object.NewEnvironment()

	env.SetWriter(os.Stdout)

	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	if f, err := os.Create(history); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	for {
		source, err := line.Prompt(">> ")

		if err == liner.ErrPromptAborted {
			fmt.Println("   Exiting...")
			os.Exit(1)
		} else {
			run(source, env)

			if errors.HadParseError {
				fmt.Printf("%s\n", errors.ParseErrorMessage)
			} else if errors.HadRuntimeError {
				fmt.Printf("%s\n", errors.RuntimeErrorMessage)
			}

			errors.HadParseError = false
			errors.HadRuntimeError = false
			line.AppendHistory(source)
		}
	}
}

func runFile(file string) {
	source, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	env := object.NewEnvironment()
	env.SetWriter(os.Stdout)

	run(string(source), env)

	if errors.HadParseError || errors.HadRuntimeError {
		os.Exit(1)
	}
}

func runServer(path string) {
	currentTime := time.Now()

	fmt.Printf("%s --> ", currentTime.Format("2006/01/02 15:04:05"))
	fmt.Printf(InfoColor, fmt.Sprintf("Starting Ghost %s server: ", version.Version))
	fmt.Printf("%s\n", flagAddress)

	http.Handle("/ghost/css/", http.StripPrefix("/ghost/css/", http.FileServer(http.Dir("./server/templates/public"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := "success"
		start := time.Now()

		env := object.NewEnvironment()
		env.SetWriter(w)

		source := readSource(path)

		run(source, env)

		if errors.HadParseError == true {
			status = "error"
			sourceBuffer, _ := ioutil.ReadFile(path)

			data := struct {
				Message string
				Host string
				Path string
				Version string
				OperatingSystem string
				Method string
				AcceptLanguage string
				AcceptEncoding string
				Source string
			}{
				Message: errors.ParseErrorMessage,
				Host: r.Host,
				Path: r.URL.Path,
				Version: version.Version,
				OperatingSystem: runtime.GOOS,
				Method: r.Method,
				AcceptLanguage: r.Header.Get("Accept-Language"),
				AcceptEncoding: r.Header.Get("Accept-Encoding"),
				Source: strings.Replace(string(sourceBuffer), "\n", "<br>", -1),
			}

			t, _ := template.ParseFiles("server/templates/error.html")
			t.Execute(w, data)

			errors.Reset()
		}

		secs := time.Since(start).String()

		log.Printf("--> %s (%s) %s (%s)", r.Method, status, r.URL.Path, secs)
	})

	log.Fatal(http.ListenAndServe(viper.GetString("SERVER_ADDRESS"), nil))
}

func readSource(path string) string {
	source, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(source)
}

func run(source string, env *object.Environment) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()

	if errors.HadParseError {
		return
	}

	interpreter.Interpret(statements, env)

	if flagTokens {
		fmt.Printf("   =====\n")

		for index, token := range tokens {
			fmt.Printf("   [%d] %s\n", index, token.String())
		}
	}
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("    ghost [flags] {file}")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("    -h  show help")
	fmt.Println("    -i  enter interactive mode after executing file")
	fmt.Println("    -v  show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("    ghost")
	fmt.Println()
	fmt.Println("            Start Ghost REPL")
	fmt.Println()
	fmt.Println("    ghost example.ghost")
	fmt.Println()
	fmt.Println("            Execute source file (example.ghost)")
	fmt.Println()
	fmt.Println("    ghost -i example.ghost")
	fmt.Println()
	fmt.Println("            Execute source file (example.ghost)")
	fmt.Println("            and enter interactive mode (REPL)")
	fmt.Println("            with the scripts environment intact")
	fmt.Println()
	fmt.Println()
}
