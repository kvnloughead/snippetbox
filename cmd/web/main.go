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

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Serve static files out of ./ui/static directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	/* Info level log statement. Arguments after the first can either be variadic, key/value pairs, or attribute pairs created by slog.String, or a similar method. */
	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(":"+*addr, mux)

	// If http.ListenAndServe returns an error, log its message and exit.
	logger.Error(err.Error())
	os.Exit(1)
}
