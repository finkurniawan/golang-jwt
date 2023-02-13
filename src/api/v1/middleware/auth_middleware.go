package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"golang-jwt/src/api/v1/helper"
	web2 "golang-jwt/src/api/v1/model/web"
	"net/http"
	"os"
	"strings"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{
		Handler: handler,
	}
}

func (a *AuthMiddleware) unauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	webResponse := web2.Response{
		Status: "OK",
		Data:   "UNAUTHORIZED",
	}

	helper.WriteToBody(w, webResponse)
}

func (a *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && (r.RequestURI == "/api/v1/user" || r.RequestURI == "/api/v1/auth") {
		a.Handler.ServeHTTP(w, r)
	} else {
		tokenAuth := r.Header.Get("Authorization")
		splitToken := strings.Split(tokenAuth, "Bearer ")

		if len(splitToken) != 2 {
			a.unauthorized(w, r)
			return
		}

		tokenAuth = splitToken[1]

		if tokenAuth == "" {
			a.unauthorized(w, r)
			return
		}

		var jwtTokenSecret = []byte(os.Getenv("JWT_TOKEN_SECRET"))
		claims := &web2.TokenClaims{}

		token, err := jwt.ParseWithClaims(tokenAuth, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtTokenSecret, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				a.unauthorized(w, r)
				return
			}
		}

		if !token.Valid {
			a.unauthorized(w, r)
			return
		}

		a.Handler.ServeHTTP(w, r)
	}
}
