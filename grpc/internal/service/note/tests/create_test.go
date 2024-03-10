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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository

	type args struct {
		ctx context.Context
		req *model.NoteInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.BeerName()
		content = gofakeit.BeerName()

		repoErr = fmt.Errorf("repo error")

		req = &model.NoteInfo{
			Title: title,
			Body:  content,
		}
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		noteRepositoryMock noteRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			noteRepositoryMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			noteRepositoryMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteRepoMock := tt.noteRepositoryMock(mc)
			service := note.NewMockService(noteRepoMock)

			newID, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
