package main

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/hiroto1220/go-playground/connect/gen/greetconnect"
	"github.com/hiroto1220/go-playground/connect/handler"
)

func main() {
	greeter := handler.NewGreetHandler()
	path, handler := greetconnect.NewGreeterHandler(greeter)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
