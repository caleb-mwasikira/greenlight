package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/gorilla/handlers"
)

func ValidateReqMethodHandler(allowedMethods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		validReqMethod := slices.Contains(allowedMethods, req.Method)

		if !validReqMethod {
			msg := fmt.Sprintf("Request method %v not allowed in route %v", req.Method, req.URL.Path)
			http.Error(res, msg, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func NewLoggingHandler(filename string) func(http.Handler) http.Handler {
	var dest io.Writer = os.Stdout

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err == nil {
		// No errors occurred during the opening of the file
		dest = file
	}

	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dest, h)
	}
}
