package main

import (
	"fmt"
	"log"
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

func main(){
	//Golang has a http.DefaultServeMux BUT for the sake of clarity, maintainablility and security, it's generally a good idea to create your own.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)
	log.Println("Starting server on: 3000")

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
	
}