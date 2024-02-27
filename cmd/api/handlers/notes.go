package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/caleb-mwasikira/greenlight/pkg/models"
)

func (app *Application) CreateNewNote(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		msg := fmt.Sprintf("Request method %v not allowed on URL %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	// Decode request body fields into a note object
	note := &models.Note{}
	err := json.NewDecoder(req.Body).Decode(note)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Implement data validation on request body fields

	id, err := app.Notes.Insert(note.Title, note.Content, note.Expires)
	if err != nil {
		msg := fmt.Sprintf("failed to create new note: %v", err)
		http.Error(res, msg, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Status-Code", "200")
	fmt.Fprintf(res, "Created new note with id %v", id)
}

func (app *Application) GetAllNotes(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		msg := fmt.Sprintf("Request method %v not allowed on URL %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	notes, err := app.Notes.GetAll(true)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	// Marshal slice of notes into []byte
	data, err := json.MarshalIndent(notes, "  ", "")
	if err != nil {
		msg := fmt.Sprintf("failed to encode notes to json data: %v", err)
		http.Error(res, msg, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

func (app *Application) GetNote(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		msg := fmt.Sprintf("Request method %v not allowed on URL %v", req.Method, req.URL)
		http.Error(res, msg, http.StatusMethodNotAllowed)
		return
	}

	paramId := req.URL.Query().Get("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		var msg string
		if len(paramId) != 0 {
			msg = "Parameter ?id must be of type number"
		} else {
			msg = "Missing parameter ?id on the request url"
		}

		http.Error(res, msg, http.StatusBadRequest)
		return
	}

	note, err := app.Notes.Get(id)
	if err != nil {
		msg := fmt.Sprintf("Note with id %v not found in the database", id)
		http.Error(res, msg, http.StatusNotFound)
		return
	}

	data, err := json.MarshalIndent(note, " ", "")
	if err != nil {
		msg := fmt.Sprintf("failed to encode note to json data: %v", err)
		http.Error(res, msg, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}
