package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const host string = "0.0.0.0"
	const port int = 8080
	addr := fmt.Sprintf("%v:%v", host, port)

	mux := http.NewServeMux()

	// mux.Handle() method expects a http.Handler() function as a 2nd argument
	// You can use the http.HandlerFunc() method of the http object to create a handler from a normal function
	// or call the mux.HandleFunc() method directly
	mux.Handle("/", http.HandlerFunc(homePage))
	mux.HandleFunc("/about", aboutPage)

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Server started on %v", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Failed to start server due to error: %v", err)
	}
}

func homePage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		msg := fmt.Sprintf("Request method %v not allowed on url %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(res, "Welcome to the home page!")
}

func aboutPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		msg := fmt.Sprintf("Request method %v not allowed on url %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(res, "This is the about page")
}
