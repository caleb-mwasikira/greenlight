package handlers

import (
	"fmt"
	"net/http"
)

func UsersPage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to the users page")
}
