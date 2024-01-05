package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/valentedev/template-server/internal/models"
)

type application struct {
	logger        *slog.Logger
	books         *models.BookModel
	templateCache map[string]*template.Template
}

func main() {

	// Flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://admin:adminpass@localhost/books?sslmode=disable", "PostgreSQL data source name")
	flag.Parse()

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// DB connection
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Instance of application struct
	app := &application{
		logger:        logger,
		books:         &models.BookModel{DB: db},
		templateCache: templateCache,
	}

	log.Printf("starting server on %s", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
