package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/robinhawiz/snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

//This application struct will hold the application-wide dependencies for the web application.
type application struct {
	logger 		   *slog.Logger
	snippets 	   *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main(){
	//Defines a new command-line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	//This parses the command-line flags, which in turn makes it possible to read in the command-line flag values and assigns it to the addr and dsn variables.
	flag.Parse()
	//A structured logger that writes to the standard out stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
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

	//Initialize a new session manager. Then we configure it to use our MySQL database as the session store, and set a lifetime of 12h.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	//Initialize application
	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db,},
		templateCache: cache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}


	srv := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		//Create a *log.Logger from our structured logger handler, which writes log entries at Error level, and assign it to the ErrorLog field.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", slog.String("addr", srv.Addr))

	err = srv.ListenAndServe()
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