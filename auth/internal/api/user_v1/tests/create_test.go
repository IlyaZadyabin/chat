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

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    int64 = 101
		name        = "John"
		email       = "john@example.com"
		role        = desc.Role_USER

		req  = &desc.CreateRequest{Info: &desc.UserInfo{Name: name, Email: email, Role: role}, Password: "secret"}
		info = &model.UserCreate{Info: &model.UserInfo{Name: name, Email: email, Password: "secret", Role: role.String()}, Password: "secret"}

		res    = &desc.CreateResponse{Id: id}
		svcErr = fmt.Errorf("svc error")
	)

	tests := []struct {
		name   string
		args   args
		want   *desc.CreateResponse
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
				m.CreateMock.Expect(ctx, info).Return(id, nil)
				return m
			},
		},
		{
			name: "service error",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  svcErr,
			mockFn: func(mc *minimock.Controller) service.UserService {
				m := serviceMocks.NewUserServiceMock(mc)
				m.CreateMock.Expect(ctx, info).Return(0, svcErr)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.mockFn(mc)
			h := api.NewUserV1Handler(svc)
			got, err := h.Create(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.ErrorContains(t, err, "failed to create user")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
