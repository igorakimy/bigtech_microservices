package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/igorakimy/bigtech_microservices/internal/model"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	repoMocks "github.com/igorakimy/bigtech_microservices/internal/repository/mocks"
	"github.com/igorakimy/bigtech_microservices/internal/service/note"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	// Запустить тест параллельно
	t.Parallel()

	type noteRepoMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")

		res = &model.Note{
			ID: id,
		}

		req = &desc.GetRequest{
			Id: id,
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         *model.Note
		err          error
		noteRepoMock noteRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteRepoMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		},
		{
			name: "service err case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			noteRepoMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteRepoMock := tt.noteRepoMock(mc)
			service := note.NewMockService(noteRepoMock)

			newNote, err := service.Get(ctx, id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newNote)
		})
	}
}
