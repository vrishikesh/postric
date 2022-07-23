package handlers

import (
	"authentication/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtSecretKey = []byte("my_secret_key")

func Signin(w http.ResponseWriter, r *http.Request) render.Renderer {
	creds := new(Credentials)

	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		return util.ErrResponse(http.StatusBadRequest, err)
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		return util.ErrResponse(http.StatusUnauthorized, err)
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return util.ErrResponse(http.StatusInternalServerError, err)
	}

	return util.Response(http.StatusAccepted, struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	})
}
