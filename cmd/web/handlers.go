package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (a *application) home (w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		a.notFound(w)
		return
	}
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	if tmpl, err := template.ParseFiles(files...); err != nil {
		a.serverError(w,r,err)
	}else{
		if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil{
			a.serverError(w,r,err)
		}
	}
	
}

func (a *application) snippetViewHandler (w http.ResponseWriter, r *http.Request){
	if id, err := strconv.Atoi(r.URL.Query().Get("id")); err != nil || id < 1 {
		a.notFound(w)
		return
	}else{
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func (a *application) snippetCreateHandler (w http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodPost){
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w,http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Creat a new snippet..."))
	
}