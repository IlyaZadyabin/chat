package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	api "chat/chat_server/internal/api/chat_v1"
	"chat/chat_server/internal/model"
	"chat/chat_server/internal/service"
	serviceMocks "chat/chat_server/internal/service/mocks"
	desc "chat/chat_server/pkg/chat_v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}
	var (
		ctx            = context.Background()
		mc             = minimock.NewController(t)
		id       int64 = 77
		req            = &desc.CreateRequest{Usernames: []string{"a", "b"}}
		modelReq       = &model.ChatCreate{Usernames: []string{"a", "b"}}
		res            = &desc.CreateResponse{Id: id}
		svcErr         = fmt.Errorf("svc error")
	)

	tests := []struct {
		name   string
		args   args
		want   *desc.CreateResponse
		err    error
		mockFn func(mc *minimock.Controller) service.ChatService
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			want: res,
			err:  nil,
			mockFn: func(mc *minimock.Controller) service.ChatService {
				m := serviceMocks.NewChatServiceMock(mc)
				m.CreateMock.Expect(ctx, modelReq).Return(id, nil)
				return m
			},
		},
		{
			name: "error",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  svcErr,
			mockFn: func(mc *minimock.Controller) service.ChatService {
				m := serviceMocks.NewChatServiceMock(mc)
				m.CreateMock.Expect(ctx, modelReq).Return(0, svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewChatV1Handler(svc)
			got, err := h.Create(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				require.ErrorContains(t, err, "failed to create chat")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}
	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		req    = &desc.DeleteRequest{Id: 10}
		svcErr = fmt.Errorf("svc error")
	)

	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(mc *minimock.Controller) service.ChatService
	}{
		{
			name:    "success",
			args:    args{ctx: ctx, req: req},
			wantErr: nil,
			mockFn: func(mc *minimock.Controller) service.ChatService {
				m := serviceMocks.NewChatServiceMock(mc)
				m.DeleteMock.Expect(ctx, req.GetId()).Return(nil)
				return m
			},
		},
		{
			name:    "error",
			args:    args{ctx: ctx, req: req},
			wantErr: svcErr,
			mockFn: func(mc *minimock.Controller) service.ChatService {
				m := serviceMocks.NewChatServiceMock(mc)
				m.DeleteMock.Expect(ctx, req.GetId()).Return(svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewChatV1Handler(svc)
			_, err := h.Delete(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
				require.ErrorContains(t, err, "failed to delete chat")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}
	var (
		ctx      = context.Background()
		mc       = minimock.NewController(t)
		req      = &desc.SendMessageRequest{From: "a", Text: "hi", Timestamp: timestamppb.New(time.Unix(0, 0).UTC())}
		modelMsg = &model.Message{From: "a", Text: "hi", Timestamp: time.Unix(0, 0).UTC()}
		svcErr   = fmt.Errorf("svc error")
	)

	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(mc *minimock.Controller) service.ChatService
	}{
		{
			name:    "success",
			args:    args{ctx: ctx, req: req},
			wantErr: nil,
			mockFn: func(mc *minimock.Controller) service.ChatService {
				m := serviceMocks.NewChatServiceMock(mc)
				m.SendMessageMock.Expect(ctx, modelMsg).Return(nil)
				return m
			},
		},
		{
			name:    "error",
			args:    args{ctx: ctx, req: req},
			wantErr: svcErr,
			mockFn: func(mc *minimock.Controller) service.ChatService {
				m := serviceMocks.NewChatServiceMock(mc)
				m.SendMessageMock.Expect(ctx, modelMsg).Return(svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewChatV1Handler(svc)
			_, err := h.SendMessage(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
				require.ErrorContains(t, err, "failed to send message")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
