package main

import (
	"fmt"
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

func (a *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData){
	ts, ok := a.templateCache[page]
	if !ok{
		err := fmt.Errorf("the template %s does not exist", page)
		a.serverError(w,r,err)
		return
	}
	w.WriteHeader(status)
	if err := ts.ExecuteTemplate(w, "base", data); err != nil{
		a.serverError(w,r,err)
	}
}