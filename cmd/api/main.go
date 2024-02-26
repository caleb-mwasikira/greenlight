package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	config "github.com/caleb-mwasikira/greenlight/cmd/api/config"
	handlers "github.com/caleb-mwasikira/greenlight/cmd/api/handlers"
	middleware "github.com/caleb-mwasikira/greenlight/cmd/api/middleware"

	// Our main.go file doesn’t actually use anything in the mysql package.
	// So if we try to import it normally the Go compiler will raise an error.
	// However, we need the driver’s init() function to run so that it can
	// register itself with the database/sql package. Thats why we import it
	// with and underscore _
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

func connectToDatabase() (*sql.DB, error) {
	var db *sql.DB

	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	// Initialise a new sql.DB object (which is not a database connection but
	// a pool of connections) based on a DSN(data source name) in the format
	// <username>:<password>@[protocol(address)]/<db-name>?[...parameters]
	// The parameter ?parseTime=True is a driver-specific parameter that informs
	// the mysql driver to convert SQL TIME and DATE fields to Go time.Time objects
	dsn := fmt.Sprintf("%v:%v@/%v?parseTime=True",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// The sql.Open function doesn’t actually create any connections, all
	// it does is initialize the pool for future use. Actual connections to the
	// database are established lazily, as and when needed for the first time.
	// So to verify that everything is set up correctly we need to use the
	// db.Ping method to create a connection and check for any errors.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	conf := parseCmdFlags()

	var (
		logDir      string = "/tmp/greenlight"
		logFilename string = path.Join(logDir, "out.log")
	)
	app := config.NewApplication(conf, "")
	mux := http.NewServeMux()

	// Connect server to database
	db, err := connectToDatabase()
	if err != nil {
		app.ErrorLog.Fatalf("failed to connect to mysql database: %v", err)
		return
	}
	app.InfoLog.Printf("connected to mysql database ...")

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize middleware
	loggingHandler := middleware.NewLoggingHandler(logFilename)
	middlewareChain := alice.New(loggingHandler)

	// mux.Handle() method expects a http.Handler() function as a 2nd argument
	// You can use the http.HandlerFunc() method of the http object to create
	// a handler from a normal function or call the mux.HandleFunc() method directly

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

	app.InfoLog.Printf("server started on %v", conf.Addr())
	err = server.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatalf("failed to start server: %v", err)
		return
	}
}
