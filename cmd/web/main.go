package main

import (
	"log"
	"net/http"
)

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