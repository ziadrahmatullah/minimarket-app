package dto

import "github.com/ziadrahmatullah/minimarket-app/entity"

type RegisterReq struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required" json:"password"`
}

func (r *RegisterReq) ToUser() *entity.User {
	return &entity.User{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginReq struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required" json:"password"`
}

func (r *LoginReq) ToUser() *entity.User {
	return &entity.User{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginRes struct {
	AccessToken string `json:"access_token"`
}
