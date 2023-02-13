package controller

import (
	"github.com/julienschmidt/httprouter"
	"golang-jwt/src/api/v1/helper"
	"golang-jwt/src/api/v1/model/web"
	"golang-jwt/src/api/v1/service"
	"net/http"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserControllerImpl(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (u *UserControllerImpl) CreateWithRefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.Header.Get("Authorization")

	response := u.UserService.CreateWithRefreshToken(r.Context(), token)

	webResponse := web.Response{
		Status: "OK",
		Data:   response,
	}

	helper.WriteToBody(w, webResponse)
}

func (u *UserControllerImpl) Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userAuthRequest := web.UserAuthRequest{}
	helper.BodyToRequest(r, &userAuthRequest)

	response := u.UserService.Auth(r.Context(), userAuthRequest)

	webResponse := web.Response{
		Status: "OK",
		Data:   response,
	}

	helper.WriteToBody(w, webResponse)
}

func (u *UserControllerImpl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userAuthRequest := web.UsersCreateRequest{}
	helper.BodyToRequest(r, &userAuthRequest)

	response := u.UserService.Create(r.Context(), userAuthRequest)

	webResponse := web.Response{
		Status: "OK",
		Data:   response,
	}

	helper.WriteToBody(w, webResponse)
}

func (u *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userUpdateRequest := web.UsersUpdateRequest{}
	helper.BodyToRequest(r, &userUpdateRequest)

	userUpdateRequest.Id = params.ByName("user_id")
	response := u.UserService.Update(r.Context(), userUpdateRequest)

	webResponse := web.Response{
		Status: "OK",
		Data:   response,
	}

	helper.WriteToBody(w, webResponse)
}

func (u *UserControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := params.ByName("user_id")
	u.UserService.Delete(r.Context(), userId)

	webResponse := web.Response{
		Status: "OK",
	}

	helper.WriteToBody(w, webResponse)
}

func (u *UserControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := params.ByName("user_id")
	response := u.UserService.FindById(r.Context(), userId)

	webResponse := web.Response{
		Status: "OK",
		Data:   response,
	}

	helper.WriteToBody(w, webResponse)
}

func (u *UserControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := u.UserService.FindAll(r.Context())

	webResponse := web.Response{
		Status: "OK",
		Data:   response,
	}

	helper.WriteToBody(w, webResponse)
}
