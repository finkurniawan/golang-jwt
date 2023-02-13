package utils

import (
	"golang-jwt/src/api/v1/model/entity"
	"golang-jwt/src/api/v1/model/web"
)

func UserResponse(user entity.Users) web.UsersResponse {
	return web.UsersResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
