package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl = "http://localhost:5000"

	createNotePrefix = "/notes"
	getNotePrefix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createNote() (Note, error) {
	noteInfo := NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.Quote(),
		Author:   gofakeit.FirstName(),
		IsPublic: gofakeit.Bool(),
	}

	b, err := json.Marshal(noteInfo)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(
		baseUrl+createNotePrefix,
		"application/json",
		bytes.NewBuffer(b),
	)
	defer func() { _ = resp.Body.Close() }()

	if err != nil {
		return Note{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return Note{}, errors.New("failed to create")
	}

	var note Note
	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func getNote(noteID int64) (Note, error) {
	var note Note

	resp, err := http.Get(baseUrl + fmt.Sprintf(getNotePrefix, noteID))
	if err != nil {
		return Note{}, nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Note{}, nil
	}

	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		return Note{}, nil
	}

	return note, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatalln("Failed to create new note", err)
	}
	log.Printf("Note created: %+v\n", note)

	note, err = getNote(note.ID)
	if err != nil {
		log.Fatalln("Error on getting note", err)
	}
	log.Printf("Note info got: %+v\n", note)
}
