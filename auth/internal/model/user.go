package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID        int64
	Info      *UserInfo
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
}

type UserInfo struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type UserCreate struct {
	Info     *UserInfo
	Password string
}

type UserUpdate struct {
	ID   int64
	Info *UserInfo
}
