package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler{
		//Golang has a http.DefaultServeMux BUT for the sake of clarity, maintainablility and security, it's generally a good idea to create your own.
		mux := http.NewServeMux()
		fileserver := http.FileServer(http.Dir("./ui/static/"))
		mux.Handle("/static/", http.StripPrefix("/static", fileserver))
		mux.HandleFunc("/", a.home)
		mux.HandleFunc("/snippet/view", a.snippetViewHandler)
		mux.HandleFunc("/snippet/create", a.snippetCreateHandler)

		//Creating a middleware chain containing our "standard" middleware which will be used for every request our application recieves.
		standard := alice.New(a.recoverPanic, a.logRequest, secureHeaders)
		return standard.Then(mux)
}
	
