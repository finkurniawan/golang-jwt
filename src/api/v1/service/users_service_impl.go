package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"golang-jwt/src/api/v1/exception"
	"golang-jwt/src/api/v1/helper"
	"golang-jwt/src/api/v1/model/entity"
	"golang-jwt/src/api/v1/model/web"
	"golang-jwt/src/api/v1/repository"
	"golang-jwt/src/api/v1/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	validate       validator.Validate
}

func (u *UserServiceImpl) Auth(ctx context.Context, request web.UserAuthRequest) web.TokenResponse {
	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	user, err := u.UserRepository.FindByEmail(
		ctx,
		tx,
		request.Email,
	)

	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	err = utils.CheckPasswordHash(user.Password, request.Password)

	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	tokenCreateRequest := &web.TokenCreateRequest{
		UserId:    user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.FirstName,
	}

	jwtExpiredTimeToken, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME_TOKEN"))
	jwtExpiredRefreshToken, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME_REFRESH_TOKEN"))

	token := web.TokenResponse{
		Token: utils.CreateToken(
			tokenCreateRequest,
			time.Duration(jwtExpiredTimeToken)),
		RefreshToken: utils.CreateToken(
			tokenCreateRequest,
			time.Duration(jwtExpiredRefreshToken)),
	}

	return token
}

func NewUserServiceImpl(userRepository repository.UserRepository, DB *sql.DB, validate validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		validate:       validate,
	}
}

func (u *UserServiceImpl) CreateWithRefreshToken(ctx context.Context, refreshToken string) web.TokenResponse {
	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	splitToken := strings.Split(refreshToken, "Bearer ")

	if len(splitToken) != 2 {
		helper.PanicError(errors.New("Create refresh token failed"))
	}

	refreshToken = splitToken[1]

	if refreshToken == "" {
		helper.PanicError(errors.New("Create refresh token failed"))
	}

	claims := utils.ClaimsToken(refreshToken)

	_, err = u.UserRepository.FindById(
		ctx,
		tx,
		claims.UserId,
	)

	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	tokenCreateRequest := &web.TokenCreateRequest{
		UserId:    claims.UserId,
		Email:     claims.Email,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
	}

	jwtExpiredTimeToken, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME_TOKEN"))
	jwtExpiredRefreshToken, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME_REFRESH_TOKEN"))

	token := web.TokenResponse{
		Token: utils.CreateToken(
			tokenCreateRequest,
			time.Duration(jwtExpiredTimeToken)),
		RefreshToken: utils.CreateToken(
			tokenCreateRequest,
			time.Duration(jwtExpiredRefreshToken)),
	}

	return token
}

func (u *UserServiceImpl) Create(ctx context.Context, request web.UsersCreateRequest) web.UsersResponse {
	err := u.validate.Struct(
		request,
	)
	helper.PanicError(err)

	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	passwordHash, err := utils.HashPassword(request.Password)
	helper.PanicError(err)

	user := entity.Users{
		Id:        utils.Uuid(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  passwordHash,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	user = *u.UserRepository.Create(
		ctx,
		tx,
		user,
	)

	return utils.UserResponse(user)
}

func (u *UserServiceImpl) Update(ctx context.Context, request web.UsersUpdateRequest) web.UsersResponse {
	err := u.validate.Struct(request)
	helper.PanicError(err)

	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	user, err := u.UserRepository.FindById(
		ctx,
		tx,
		request.Id,
	)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.UpdatedAt = time.Now().Unix()

	user = u.UserRepository.Update(
		ctx,
		tx,
		*user,
	)

	return utils.UserResponse(*user)
}

func (u *UserServiceImpl) Delete(ctx context.Context, userId string) {
	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	user, err := u.UserRepository.FindById(ctx,
		tx, userId)
	helper.PanicError(err)

	u.UserRepository.Delete(
		ctx,
		tx,
		*user,
	)
}

func (u *UserServiceImpl) FindById(ctx context.Context, userId string) web.UsersResponse {
	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	user, err := u.UserRepository.FindById(
		ctx,
		tx,
		userId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	return utils.UserResponse(*user)
}

func (u *UserServiceImpl) FindAll(ctx context.Context) []web.UsersResponse {
	tx, err := u.DB.Begin()
	helper.PanicError(err)
	defer helper.Defer(tx)

	users := u.UserRepository.FindAll(
		ctx, tx)

	var userResponse []web.UsersResponse

	for _, user := range users {
		userResponse = append(userResponse, utils.UserResponse(user))
	}

	return userResponse
}
