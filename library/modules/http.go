package modules

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var HttpMethods = map[string]*object.LibraryFunction{}
var HttpProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(HttpMethods, "handle", httpHandle)
	RegisterMethod(HttpMethods, "listen", httpListen)
}

func httpHandle(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("http.handle", tok, args, 2); err != nil {
		return err
	}

	path, err := stringAt("http.handle", tok, args, 0)

	if err != nil {
		return err
	}

	callback, err := functionAt("http.handle", tok, args, 1)

	if err != nil {
		return err
	}

	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		// A request handler runs on its own goroutine, outside the recovery
		// that wraps a script's main run. A panic here would take the whole
		// server down with a Go traceback, so it is caught where it happens and
		// answered with a 500.
		defer recoverHandler(writer, path)

		scope.Environment.SetWriter(writer)

		requestBodyBuf := new(bytes.Buffer)
		requestBodyBuf.ReadFrom(request.Body)

		httpRequest := object.NewMap(map[string]interface{}{
			"method":        request.Method,
			"host":          request.Host,
			"contentLength": request.ContentLength,
			"protocol":      request.Proto,
			"protocolMajor": request.ProtoMajor,
			"protocolMinor": request.ProtoMinor,
			"body":          requestBodyBuf.String(),
		})

		callback.Evaluate([]object.Object{httpRequest}, writer)
	})

	return nil
}

// recoverHandler turns a panic inside a request handler into a 500 and a line
// on the server's log, rather than the end of the process.
func recoverHandler(writer http.ResponseWriter, path string) {
	recovered := recover()

	if recovered == nil {
		return
	}

	log.Error("the handler for %s stopped unexpectedly: %v", path, recovered)

	writer.WriteHeader(http.StatusInternalServerError)
}

func httpListen(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("http.listen", tok, args, 1, 2); err != nil {
		return err
	}

	port, err := integerAt("http.listen", tok, args, 0)

	if err != nil {
		return err
	}

	if port < 1 || port > 65535 {
		return object.NewError(fault.Value, tok, "`http.listen()` expects a port between 1 and 65535, got %d", port)
	}

	var ready *object.Function

	if len(args) == 2 {
		ready, err = functionAt("http.listen", tok, args, 1)

		if err != nil {
			return err
		}
	}

	address := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: address}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)

		if failure := server.Shutdown(ctx); failure != nil {
			log.Debug("could not shut the server down cleanly: %v", failure)
		}

		close(done)
	}()

	if ready != nil {
		ready.Evaluate(nil, nil)
	}

	// A server that cannot bind its port has not started, and reporting that as
	// a debug line leaves the script believing it is listening. It is a failure
	// of the world outside the program, so it comes back as one.
	if failure := server.ListenAndServe(); failure != nil && failure != http.ErrServerClosed {
		return object.NewError(fault.System, tok, "`http.listen()` could not listen on port %d: %s", port, failure)
	}

	<-done

	return nil
}
