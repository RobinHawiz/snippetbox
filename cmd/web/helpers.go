package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error){
	a.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()), slog.String("trace", string(debug.Stack())))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *application) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}

func (a *application) notFound(w http.ResponseWriter){
	a.clientError(w, http.StatusNotFound)
}
