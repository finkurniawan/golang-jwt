package exception

import (
	"golang-jwt/src/api/v1/helper"
	"golang-jwt/src/api/v1/model/web"
	"net/http"
	"strconv"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if unauthorized(w, r, err) {
		return
	} else if badRequest(w, r, err) {
		return
	} else if notFound(w, r, err) {
		return
	}
	internalServerError(w, r, err)
}

func badRequest(w http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		response := web.Response{
			Status: strconv.Itoa(http.StatusNotFound),
			Data:   exception.Error,
		}

		helper.WriteToBody(w, response)
		return true
	} else {
		return false
	}
}

func internalServerError(w http.ResponseWriter, _ *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := web.Response{
		Status: strconv.Itoa(http.StatusInternalServerError),
		Data:   err,
	}

	helper.WriteToBody(w, response)
}

func notFound(w http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		response := web.Response{
			Status: strconv.Itoa(http.StatusNotFound),
			Data:   exception.Error,
		}

		helper.WriteToBody(w, response)
		return true
	} else {
		return false
	}
}

func unauthorized(w http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(UnauthorizedError)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := web.Response{
			Status: strconv.Itoa(http.StatusUnauthorized),
			Data:   exception.Error,
		}

		helper.WriteToBody(w, response)
		return true
	} else {
		return false
	}
}
