package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// A struct containing application-wide dependencies.
type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", "4000", "HTTP Network Address")
	flag.Parse()

	// Initialize structured logger to stdout with default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // include file and line number
	}))

	app := &application{logger: logger}

	/* Info level log statement. Arguments after the first can either be variadic, key/value pairs, or attribute pairs created by slog.String, or a similar method. */
	logger.Info("starting server", slog.String("addr", *addr))

	mux := app.routes()
	err := http.ListenAndServe(":"+*addr, mux)

	// If http.ListenAndServe returns an error, log its message and exit.
	logger.Error(err.Error())
	os.Exit(1)
}
