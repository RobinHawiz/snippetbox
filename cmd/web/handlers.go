package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
) 

type snippet struct {
	content string
}

func home (w http.ResponseWriter, r *http.Request){
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}else{
		if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil{
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	
}

func snippetViewHandler (w http.ResponseWriter, r *http.Request){
	if id, err := strconv.Atoi(r.URL.Query().Get("id")); err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}else{
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func snippetCreateHandler (w http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodPost){
		w.Header().Set("Allow", http.MethodPost)
		//Calls the w.WriteHeader() and w.Write() methods for us. Pretty neat.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Creat a new snippet..."))
	
}