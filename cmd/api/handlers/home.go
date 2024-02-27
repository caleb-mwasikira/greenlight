package handlers

import (
	"fmt"
	"net/http"
)

func (app *Application) HomePage(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	if req.Method != "GET" {
		msg := fmt.Sprintf("Request method %v not allowed on URL %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(res, "Welcome to the home page!")
}

func (app *Application) AboutPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		msg := fmt.Sprintf("Request method %v not allowed on URL %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(res, "This is the about page")
}
