package handlers

import (
	"fmt"
	"net/http"
)

func HomePage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Welcome to the home page!")
}

func AboutPage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "This is the about page")
}
