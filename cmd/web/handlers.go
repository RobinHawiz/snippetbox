package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
)

func (a *application) home (w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	if tmpl, err := template.ParseFiles(files...); err != nil {
		a.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}else{
		if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil{
			a.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	
}

func (a *application) snippetViewHandler (w http.ResponseWriter, r *http.Request){
	if id, err := strconv.Atoi(r.URL.Query().Get("id")); err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}else{
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func (a *application) snippetCreateHandler (w http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodPost){
		w.Header().Set("Allow", http.MethodPost)
		//Calls the w.WriteHeader() and w.Write() methods for us. Pretty neat.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Creat a new snippet..."))
	
}