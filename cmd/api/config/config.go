package config

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Config struct {
	Host      string // IP address of host machine
	Port      int
	StaticDir string // Path to static assets
}

func (config *Config) Addr() string {
	return net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
}

type Application struct {
	Config   *Config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewApplication(config *Config, logFile string) *Application {
	var dest io.Writer = os.Stdout

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err == nil {
		// No errors occurred during the opening of the file
		dest = file
	}

	return &Application{
		Config:   config,
		InfoLog:  log.New(dest, "INFO\t", log.LstdFlags|log.Lshortfile),
		ErrorLog: log.New(dest, "ERROR\t", log.LstdFlags|log.Lshortfile),
	}
}
