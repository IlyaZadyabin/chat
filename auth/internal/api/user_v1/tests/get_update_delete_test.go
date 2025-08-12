package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	api "chat/auth/internal/api/user_v1"
	"chat/auth/internal/model"
	"chat/auth/internal/service"
	serviceMocks "chat/auth/internal/service/mocks"
	desc "chat/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}
	var (
		ctx          = context.Background()
		mc           = minimock.NewController(t)
		id     int64 = 42
		req          = &desc.GetRequest{Id: id}
		user         = &model.User{ID: id, Info: &model.UserInfo{Name: "Jane", Email: "jane@example.com", Role: "USER"}}
		res          = &desc.GetResponse{User: &desc.User{Id: id, Info: &desc.UserInfo{Name: "Jane", Email: "jane@example.com", Role: desc.Role_USER}}}
		svcErr       = fmt.Errorf("svc error")
	)

	tests := []struct {
		name   string
		args   args
		want   *desc.GetResponse
		err    error
		mockFn func(mc *minimock.Controller) service.UserService
	}{
		{
			name: "success",
			args: args{ctx: ctx, req: req},
			want: res,
			err:  nil,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.GetMock.Expect(ctx, id).Return(user, nil)
				return m
			},
		},
		{
			name: "error",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  svcErr,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.GetMock.Expect(ctx, id).Return(nil, svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewUserV1Handler(svc)
			got, err := h.Get(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.ErrorContains(t, err, "failed to get user")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}
	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		req    = &desc.UpdateRequest{Id: 7, Info: &desc.UpdateUserInfo{}}
		svcErr = fmt.Errorf("svc error")
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(mc *minimock.Controller) service.UserService
	}{
		{
			name:    "success",
			args:    args{ctx: ctx, req: req},
			wantErr: nil,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.UpdateMock.Expect(ctx, &model.UserUpdate{ID: req.GetId(), Info: &model.UserInfo{Role: "USER"}}).Return(nil)
				return m
			},
		},
		{
			name:    "error",
			args:    args{ctx: ctx, req: req},
			wantErr: svcErr,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.UpdateMock.Expect(ctx, &model.UserUpdate{ID: req.GetId(), Info: &model.UserInfo{Role: "USER"}}).Return(svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewUserV1Handler(svc)
			_, err := h.Update(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.ErrorContains(t, err, "failed to update user")
			} else {
				require.NoError(t, err)
			}
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
		req    = &desc.DeleteRequest{Id: 11}
		svcErr = fmt.Errorf("svc error")
	)

	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(mc *minimock.Controller) service.UserService
	}{
		{
			name:    "success",
			args:    args{ctx: ctx, req: req},
			wantErr: nil,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.DeleteMock.Expect(ctx, req.GetId()).Return(nil)
				return m
			},
		},
		{
			name:    "error",
			args:    args{ctx: ctx, req: req},
			wantErr: svcErr,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.DeleteMock.Expect(ctx, req.GetId()).Return(svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewUserV1Handler(svc)
			_, err := h.Delete(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.ErrorContains(t, err, "failed to delete user")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
