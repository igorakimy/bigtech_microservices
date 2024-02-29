package converter

import (
	"github.com/igorakimy/bigtech_microservices/internal/repository/note/model"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromRepo(note *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromRepo(&note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNotesFromRepo(notes []model.Note) []*desc.Note {
	var descNotes []*desc.Note

	for _, note := range notes {
		n := &desc.Note{
			Id: note.ID,
			Info: &desc.NoteInfo{
				Title:   note.Info.Title,
				Content: note.Info.Body,
			},
		}

		descNotes = append(descNotes, n)
	}

	return descNotes
}

func ToNoteInfoFromRepo(info *model.Info) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Body,
	}
}
