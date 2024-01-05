package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/kvnloughead/snippetbox/internal/models"

	// Aliasing with a blank identifier because the driver isn't used explicitly.
	_ "github.com/go-sql-driver/mysql"
)

// A struct containing application-wide dependencies.
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", "4000", "HTTP Network Address")
	dsn := flag.String(
		"dsn",
		"web:devpass@/snippetbox?parseTime=true",
		"MySQL data source name (aka 'connection string')")
	flag.Parse()

	// Initialize structured logger to stdout with default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // include file and line number
	}))

	// Initialize sql.DB connection pool for the provided DSN.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Initialize template cache and snippet modeland add to app struct.
	snippets := &models.SnippetModel{DB: db}
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      snippets,
		templateCache: templateCache,
	}

	/* Info level log statement. Arguments after the first can either be variadic, key/value pairs, or attribute pairs created by slog.String, or a similar method. */
	logger.Info("starting server", slog.String("addr", *addr))

	mux := app.routes()
	err = http.ListenAndServe(":"+*addr, mux)

	// If http.ListenAndServe returns an error, log its message and exit.
	logger.Error(err.Error())
	os.Exit(1)
}

// Returns an sql.DB connection pool for the supplied data source name (DSN).
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify that the connection is alive, reestablishing it if necessary.
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
