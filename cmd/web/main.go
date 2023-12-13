package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

//This application struct will hold the application-wide dependencies for the web application.
type application struct {
	logger *slog.Logger
}

func main(){
	//Defines a new command-line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")
	//This parses the command-line flag, which in turn makes it possible to read in the command-line flag value and assigns it to the addr variable.
	//Ex: go run . -addr=":9999"
	flag.Parse()
	//A structured logger that writes to the standard out stream and uses customized settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}))
	//Initialize application
	app := &application{
		logger: logger,
	}
	//Golang has a http.DefaultServeMux BUT for the sake of clarity, maintainablility and security, it's generally a good idea to create your own.
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetViewHandler)
	mux.HandleFunc("/snippet/create", app.snippetCreateHandler)

	logger.Info("Starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
	
}