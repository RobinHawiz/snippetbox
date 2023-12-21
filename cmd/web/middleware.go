package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w,r)
	})
}

func (a *application) logRequest(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		var (
			ip = r.RemoteAddr
			proto = r.Proto
			method = r.Method
			uri = r.URL.RequestURI()
		)

		a.logger.Info("recieved request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w,r)
	})
}

func (a *application) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w,r,fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w,r)
	})
}

func (a *application) requireAuthentication(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//Redirect unauthorized users to the login page.
		if !(a.isAuthenticated(r)) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}

		//Pages that require authentication can't be stored in the users browser cache (or other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w,r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})
	return csrfHandler
}

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//Checks if the user is logged in. Or in other words, if the current session has a authenticatedUserID. If not, it returns 0 and we call the next handler in the chain.
		id := a.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w,r)
			return
		}

		//Otherwise, we check to see it a user with that ID exists in our database.
		exists, err := a.users.Exists(id)
		if err != nil {
			a.serverError(w,r,err)
			return
		}

		//If a matching user is found, we know that the request is coming from an authenticated user who exists in our dabatase.
		if exists {
			//We create a copy of the request with a key and true value stored in the request context. We then pass this copy to the next handler in the chain.
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w,r)
	})
}