package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	const host string = "0.0.0.0"
	const port int = 8080
	addr := fmt.Sprintf("%v:%v", host, port)

	router := httprouter.New()

	router.GET("/", homePage)
	router.GET("/about", aboutPage)
	router.ServeFiles("/static/*filepath", http.Dir("./public"))

	log.Printf("Starting server on %v", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Failed to start server due to error: %v", err)
	}
}

func homePage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprint(res, "Welcome to the home page!")
}

func aboutPage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprintf(res, "This is the about page")
}
