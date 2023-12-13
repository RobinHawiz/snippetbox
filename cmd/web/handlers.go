package main

import (
	"fmt"
	"net/http"
	"strconv"
) 

type snippet struct {
	content string
}

func home (w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello World"))
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