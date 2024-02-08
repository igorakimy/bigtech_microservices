package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	baseUrl = ":5000"

	createPostfix = "/notes"
	getPostfix    = "/notes/{id}"
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

type SyncMap struct {
	elems map[int64]*Note
	mu    sync.RWMutex
}

var notes = &SyncMap{
	elems: map[int64]*Note{},
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &NoteInfo{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "Failed to decode note data", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	randID := rand.Int63()

	note := &Note{
		ID:        randID,
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}

	notes.mu.Lock()
	defer notes.mu.Unlock()

	notes.elems[note.ID] = note
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")
	id, err := parseNoteID(noteID)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	notes.mu.RLock()
	defer notes.mu.RUnlock()

	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}
}

func parseNoteID(noteID string) (int64, error) {
	id, err := strconv.ParseInt(noteID, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func main() {
	router := chi.NewRouter()

	router.Get(getPostfix, getNoteHandler)
	router.Post(createPostfix, createNoteHandler)

	log.Fatal(http.ListenAndServe(baseUrl, router))
}
