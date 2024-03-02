package converter

import (
	"github.com/igorakimy/bigtech_microservices/internal/model"
	modelRepo "github.com/igorakimy/bigtech_microservices/internal/repository/note/model"
)

func ToNoteFromRepo(note *modelRepo.Note) *model.Note {
	return &model.Note{
		ID:        note.ID,
		Info:      ToNoteInfoFromRepo(&note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNotesFromRepo(notes []modelRepo.Note) []model.Note {
	var modelNotes []model.Note

	for _, note := range notes {
		n := model.Note{
			ID: note.ID,
			Info: model.NoteInfo{
				Title: note.Info.Title,
				Body:  note.Info.Body,
			},
		}

		modelNotes = append(modelNotes, n)
	}

	return modelNotes
}

func ToNoteInfoFromRepo(info *modelRepo.Info) model.NoteInfo {
	return model.NoteInfo{
		Title: info.Title,
		Body:  info.Body,
	}
}
