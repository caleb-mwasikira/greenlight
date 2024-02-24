package main

import (
	"flag"
	"net/http"
	"path"

	config "github.com/caleb-mwasikira/greenlight/cmd/api/config"
	handlers "github.com/caleb-mwasikira/greenlight/cmd/api/handlers"
	middleware "github.com/caleb-mwasikira/greenlight/cmd/api/middleware"

	"github.com/justinas/alice"
)

func parseCmdFlags() *config.Config {
	conf := &config.Config{}

	flag.StringVar(&conf.Host, "host", "127.0.0.1", "HTTP network address")
	flag.IntVar(&conf.Port, "port", 8080, "Port number to run the web server")
	flag.StringVar(&conf.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	return conf
}

func main() {
	conf := parseCmdFlags()

	var (
		logDir      string = "/tmp/greenlight"
		logFilename string = path.Join(logDir, "out.log")
	)
	app := config.NewApplication(conf, "")
	mux := http.NewServeMux()

	// Initialize middleware
	loggingHandler := middleware.NewLoggingHandler(logFilename)
	middlewareChain := alice.New(loggingHandler)

	// mux.Handle() method expects a http.Handler() function as a 2nd argument
	// You can use the http.HandlerFunc() method of the http object to create a handler from a normal function
	// or call the mux.HandleFunc() method directly

	// Example: mux.Handle("/", http.HandlerFunc(app.HomePage))
	mux.Handle("/", middlewareChain.ThenFunc(handlers.HomePage))
	mux.Handle("/about", middlewareChain.ThenFunc(handlers.AboutPage))
	mux.Handle("/users", middlewareChain.ThenFunc(handlers.UsersPage))

	fileServer := http.FileServer(http.Dir(conf.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	server := &http.Server{
		Addr:     conf.Addr(),
		Handler:  mux,
		ErrorLog: app.ErrorLog,
	}

	app.InfoLog.Printf("Server started on %v", conf.Addr())
	err := server.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatalf("Failed to start server due to error: %v", err)
	}
}
