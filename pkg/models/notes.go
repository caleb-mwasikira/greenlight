package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrNoRecord = errors.New("no matching record found for model")

type Note struct {
	Id      int       `json:"id,omitempty"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

type NoteRepository struct {
	DB *sql.DB
}

func (m *NoteRepository) Insert(title, content string, expires time.Time) (int, error) {
	createdAt := time.Now().UTC().Format("2006-01-02 15:04:05")
	expiresAt := expires.UTC().Format("2006-01-02 15:04:05")

	stmt := `INSERT INTO notes (title, content, created, expires)
	VALUES(?, ?, ?, ?);`

	result, err := m.DB.Exec(stmt, title, content, createdAt, expiresAt)
	if err != nil {
		return 0, err
	}

	// Use lastInsertedID() method of the result object to get the
	// ID of our newly inserted record in the notes table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *NoteRepository) Get(id int) (*Note, error) {
	stmt := `SELECT * FROM notes WHERE id = ? AND expires > UTC_TIMESTAMP();`

	// *sql.DB.QueryRow() function returns a sql.Row object which holds
	// results from the database
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new Note object
	note := &Note{}

	err := row.Scan(&note.Id, &note.Title, &note.Content, &note.Created, &note.Expires)
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("note with id of %v was not found in the database", id)
		return nil, errors.New(msg)
	} else if err != nil {
		return nil, err
	}

	// Everything went OK
	return note, nil
}

func (m *NoteRepository) GetAll(includeExpiredNotes bool) ([]*Note, error) {
	var stmt string
	if includeExpiredNotes {
		stmt = `SELECT * FROM notes;`
	} else {
		stmt = `SELECT * FROM notes expires > UTC_TIMESTAMP();`
	}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		msg := fmt.Sprintf("failed to fetch notes from database: %v", err)
		return nil, errors.New(msg)
	}
	defer rows.Close()

	// Initialize an slice of notes
	notes := []*Note{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by
	// the rows.Scan() method. If iteration over all the rows completes
	// then the resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		note := &Note{}

		err := rows.Scan(&note.Id, &note.Title, &note.Content, &note.Created, &note.Expires)
		if err != nil {
			return nil, err
		}

		// Append the new note to the slice of notes
		notes = append(notes, note)
	}

	// We call rows.Err() to retrieve any errors encoutered during
	// record iteration by the rows.Next() method
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Everything is OK
	return notes, nil
}
