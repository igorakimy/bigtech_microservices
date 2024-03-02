package converter

import (
	"database/sql"
	"github.com/igorakimy/bigtech_microservices/internal/model"
	modelRepo "github.com/igorakimy/bigtech_microservices/internal/repository/note/model"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromDesc(note *desc.Note) *model.Note {
	return &model.Note{
		ID:        note.Id,
		Info:      *ToNoteInfoFromDesc(note.Info),
		CreatedAt: note.CreatedAt.AsTime(),
		UpdatedAt: sql.NullTime{
			Time: note.UpdatedAt.AsTime(),
		},
	}
}

func ToNotesFromDesc(notes []modelRepo.Note) []model.Note {
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

func ToNoteInfoFromDesc(info *desc.NoteInfo) *model.NoteInfo {
	return &model.NoteInfo{
		Title: info.Title,
		Body:  info.Content,
	}
}

func ToUpdateNoteInfoFromDesc(info *desc.UpdateNoteInfo) *model.UpdateNoteInfo {
	return &model.UpdateNoteInfo{
		Title: info.GetTitle().GetValue(),
		Body:  info.GetContent().GetValue(),
	}
}

func ToNoteFromService(note *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromService(note.Info),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNotesFromService(notes []model.Note) []*desc.Note {
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

func ToNoteInfoFromService(info model.NoteInfo) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Body,
	}
}
