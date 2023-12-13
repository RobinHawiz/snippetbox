package main

import (
	"flag"
	"log"
	"net/http"
)

func main(){
	//Defines a new command-line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")
	//This parses the command-line flag, which in turn makes it possible to read in the command-line flag value and assigns it to the addr variable.
	//Ex: go run . -addr=":9999"
	flag.Parse()
	//Golang has a http.DefaultServeMux BUT for the sake of clarity, maintainablility and security, it's generally a good idea to create your own.
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)
	log.Printf("Starting server on %s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
	
}