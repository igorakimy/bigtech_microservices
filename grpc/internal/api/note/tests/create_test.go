package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/igorakimy/bigtech_microservices/internal/api/note"
	"github.com/igorakimy/bigtech_microservices/internal/model"
	"github.com/igorakimy/bigtech_microservices/internal/service"
	serviceMocks "github.com/igorakimy/bigtech_microservices/internal/service/mocks"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.BeerName()
		content = gofakeit.BeerName()

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.NoteInfo{
				Title:   title,
				Content: content,
			},
		}

		info = &model.NoteInfo{
			Title: title,
			Body:  content,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		noteServiceMock noteServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			noteServiceMock := tt.noteServiceMock(mc)
			api := note.NewServerImplementation(noteServiceMock)

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
