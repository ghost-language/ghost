package modules

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
)

var Http = map[string]*object.LibraryFunction{}

func init() {
	RegisterMethod(Http, "handle", httpHandle)
	RegisterMethod(Http, "listen", httpListen)
}

func httpHandle(args ...object.Object) object.Object {
	// path
	if args[0].Type() != object.STRING {
		return nil
	}

	// callback
	if args[1].Type() != object.FUNCTION {
		return nil
	}

	path := args[0].(*object.String).Value

	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})

	return nil
}

func httpListen(args ...object.Object) object.Object {
	if args[0].Type() != object.NUMBER {
		return nil
	}

	if len(args) == 2 {
		if args[1].Type() != object.FUNCTION {
			return nil
		}
	}

	port := args[0].(*object.Number).String()

	server := &http.Server{
		Addr: ":" + port,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)

		if err := server.Shutdown(ctx); err != nil {
			log.Debug("Could not gracefull shutdown the server: %v\n", err)
		}

		close(done)
	}()

	if len(args) == 2 {
		callback := args[1].(*object.Function)

		callback.Evaluate(nil)
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Debug("Could not listen on %s: %v", port, err)
	}

	<-done

	return nil
}
