package service

import (
	"context"
	web2 "golang-jwt/src/api/v1/model/web"
)

type UserService interface {
	Create(ctx context.Context, request web2.UsersCreateRequest) web2.UsersResponse
	Update(ctx context.Context, request web2.UsersUpdateRequest) web2.UsersResponse
	Delete(ctx context.Context, userId string)
	FindById(ctx context.Context, userId string) web2.UsersResponse
	FindAll(ctx context.Context) []web2.UsersResponse
	Auth(ctx context.Context, request web2.UserAuthRequest) web2.TokenResponse
	CreateWithRefreshToken(ctx context.Context, refreshToken string) web2.TokenResponse
}
