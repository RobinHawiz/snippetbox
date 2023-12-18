package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/robinhawiz/snippetbox/internal/models"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

//This application struct will hold the application-wide dependencies for the web application.
type application struct {
	logger 		  *slog.Logger
	snippets 	  *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main(){
	//Defines a new command-line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	//This parses the command-line flag, which in turn makes it possible to read in the command-line flag value and assigns it to the addr variable.
	flag.Parse()
	//A structured logger that writes to the standard out stream and uses customized settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	cache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//Initialize a decoder instance.
	formDecoder := form.NewDecoder()

	//Initialize application
	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db,},
		templateCache: cache,
		formDecoder: formDecoder,
	}

	logger.Info("Starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
	
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}