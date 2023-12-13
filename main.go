package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreateHandler (w http.ResponseWriter, r *http.Request){
	if(r.Method != "POST"){
		w.Header().Set("Allow", "POST")
		//Calls the w.WriteHeader() and w.Write() methods for us. Pretty neat.
		http.Error(w, "Method Not Allowed", 405)
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