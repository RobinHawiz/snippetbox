package main

import "net/http"

func (a *application) routes() *http.ServeMux{
		//Golang has a http.DefaultServeMux BUT for the sake of clarity, maintainablility and security, it's generally a good idea to create your own.
		mux := http.NewServeMux()
		fileserver := http.FileServer(http.Dir("./ui/static/"))
		mux.Handle("/static/", http.StripPrefix("/static", fileserver))
		mux.HandleFunc("/", a.home)
		mux.HandleFunc("/snippet/view", a.snippetViewHandler)
		mux.HandleFunc("/snippet/create", a.snippetCreateHandler)
		return mux
}
	
