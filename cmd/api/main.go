package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Host      string // IP address of host machine
	Port      int
	StaticDir string // Path to static assets
}

func (cfg *Config) Addr() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}

func parseCmdFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Host, "host", "127.0.0.1", "HTTP network address")
	flag.IntVar(&cfg.Port, "port", 8080, "Port number to run the web server")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	return cfg
}

func main() {
	cfg := parseCmdFlags()
	mux := http.NewServeMux()

	// mux.Handle() method expects a http.Handler() function as a 2nd argument
	// You can use the http.HandlerFunc() method of the http object to create a handler from a normal function
	// or call the mux.HandleFunc() method directly
	mux.Handle("/", http.HandlerFunc(homePage))
	mux.HandleFunc("/about", aboutPage)

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Create a logger for writing information and error messages
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.LstdFlags|log.Lshortfile)

	server := &http.Server{
		Addr:     cfg.Addr(),
		Handler:  mux,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Server started on %v", cfg.Addr())
	err := server.ListenAndServe()
	if err != nil {
		errorLog.Fatalf("Failed to start server due to error: %v", err)
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
