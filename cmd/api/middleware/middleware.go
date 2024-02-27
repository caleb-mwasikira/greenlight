package middleware

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

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
