package main

import (
	"crypto/tls"
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
	users		   *models.UserModel
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
	//Session cookie will only be sent by a user's web browser when a HTTPS connection is being used.
	sessionManager.Cookie.Secure = true

	//Initialize application
	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db,},
		users: &models.UserModel{DB: db,},
		templateCache: cache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	//Initialize tls.Config struct to hold non-default TLS settings thath we want the server to use.
	tlsConfig := &tls.Config{
		//We restrict the elliptic curves that can potentially be used during the TLS handshake.
		CurvePreferences: []tls.CurveID{tls.X25519},
		//We only want to support cipher suites which use ECDHE and not support weaker suites that use RC4, 3DES, CBC.
		//Note that if a TLS 1.3 connection is negotiated, any CipherSuites field will be ignored (The suites Go supports for TLS 1.3 connections are considered to be safe).
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
	}

	//Initialize server
	srv := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		//Create a *log.Logger from our structured logger handler, which writes log entries at Error level, and assign it to the ErrorLog field.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		//Add idle, Read and Write timeouts to the server.
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", slog.String("addr", srv.Addr))

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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