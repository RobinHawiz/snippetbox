package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/robinhawiz/snippetbox/ui"
)

func (a *application) routes() http.Handler{
		//Initialize router.
		router := httprouter.New()
		//We set our own helper function to be called when the router needs to send a 404 response. Router will otherwise by default use http.NotFound().
		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			a.notFound(w)
		})

		//Load the static files into the website.
		fileserver := http.FileServer(http.FS(ui.Files))
		router.Handler("GET", "/static/*filepath", fileserver)

		//Unprotected application routes using the "dynamic" middleware chain.
		dynamic := alice.New(a.sessionManager.LoadAndSave, noSurf, a.authenticate)
		
		router.Handler("GET", "/", dynamic.ThenFunc(a.home))
		router.Handler("GET", "/snippet/view/:id", dynamic.ThenFunc(a.snippetViewHandler))
		router.Handler("GET", "/user/signup", dynamic.ThenFunc(a.userSignup))
		router.Handler("POST", "/user/signup", dynamic.ThenFunc(a.userSignupPost))
		router.Handler("GET", "/user/login", dynamic.ThenFunc(a.userLogin))
		router.Handler("POST", "/user/login", dynamic.ThenFunc(a.userLoginPost))


		//Protected (authenticated-only) application routes, using a new "protected" middleware chain which includes the requireAuthentication middleware.
		protected := dynamic.Append(a.requireAuthentication)

		router.Handler("GET", "/snippet/create", protected.ThenFunc(a.snippetCreateHandler))
		router.Handler("POST", "/snippet/create", protected.ThenFunc(a.snippetCreatePostHandler))
		router.Handler("POST", "/user/logout", protected.ThenFunc(a.userLogoutPost))

		//Creating a middleware chain containing our "standard" middleware which will be used for every request our application recieves.
		standard := alice.New(a.recoverPanic, a.logRequest, secureHeaders)
		return standard.Then(router)
}
	
