package converter

import (
	"chat/auth/internal/model"
	desc "chat/auth/pkg/user_v1"
)

func ToUserCreateFromDesc(user *desc.CreateRequest) *model.UserCreate {
	return &model.UserCreate{
		Info: &model.UserInfo{
			Name:     user.Info.Name,
			Email:    user.Info.Email,
			Password: user.Password,
			Role:     user.Info.Role.String(),
		},
		Password: user.Password,
	}
}

func ToUserUpdateFromDesc(user *desc.UpdateRequest) *model.UserUpdate {
	return &model.UserUpdate{
		ID: user.Id,
		Info: &model.UserInfo{
			Name:     user.Info.Name.GetValue(),
			Email:    user.Info.Email.GetValue(),
			Password: "", // Password is not updated here
			Role:     user.Info.Role.String(),
		},
	}
}

func ToUserFromService(user *model.User) *desc.User {
	return &desc.User{
		Id: user.ID,
		Info: &desc.UserInfo{
			Name:  user.Info.Name,
			Email: user.Info.Email,
			Role:  desc.Role(desc.Role_value[user.Info.Role]),
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
