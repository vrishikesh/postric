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
		return util.ErrBadRequest(err)
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		return util.ErrUnauthorized(err)
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
		return util.ErrInternalServer(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return util.Response(struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	})
}
