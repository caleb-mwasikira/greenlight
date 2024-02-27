package handlers

import (
	"database/sql"
	"io"
	"log"
	"os"

	config "github.com/caleb-mwasikira/greenlight/cmd/api/config"
	models "github.com/caleb-mwasikira/greenlight/pkg/models"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	// Repositories of where our application fetches data from
	// (currently fetching from the database)
	Notes *models.NoteRepository
}

func NewApplication(conf *config.Config, db *sql.DB) *Application {
	var dest io.Writer = os.Stdout

	file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		dest = os.Stdout
	} else {
		dest = file
	}

	return &Application{
		InfoLog:  log.New(dest, "INFO\t", log.LstdFlags|log.Lshortfile),
		ErrorLog: log.New(dest, "ERROR\t", log.LstdFlags|log.Lshortfile),
		Notes: &models.NoteRepository{
			DB: db,
		},
	}
}
