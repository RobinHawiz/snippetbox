package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler{
		//Initialize router.
		router := httprouter.New()
		//We set our own helper function to be called when the router needs to send a 404 response. Router will otherwise by default use http.NotFound().
		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			a.notFound(w)
		})

		//Load the static files into the website
		fileserver := http.FileServer(http.Dir("./ui/static/"))
		router.Handler("GET", "/static/*filepath", http.StripPrefix("/static", fileserver))

		//Routing
		router.HandlerFunc("GET", "/", a.home)
		router.HandlerFunc("GET", "/snippet/view/:id", a.snippetViewHandler)
		router.HandlerFunc("GET", "/snippet/create", a.snippetCreateHandler)
		router.HandlerFunc("POST", "/snippet/create", a.snippetCreatePostHandler)

		//Creating a middleware chain containing our "standard" middleware which will be used for every request our application recieves.
		standard := alice.New(a.recoverPanic, a.logRequest, secureHeaders)
		return standard.Then(router)
}
	
