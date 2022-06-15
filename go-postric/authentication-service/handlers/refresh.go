package handlers

import (
	"authentication/util"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
)

func Refresh(w http.ResponseWriter, r *http.Request) render.Renderer {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return util.ErrUnauthorized(err)
		}

		return util.ErrBadRequest(err)
	}

	tknStr := c.Value
	claims := new(Claims)
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return util.ErrUnauthorized(err)
		}

		return util.ErrBadRequest(err)
	}
	if !tkn.Valid {
		return util.ErrUnauthorized(err)
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 30*time.Second {
		return util.ErrBadRequest(err)
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
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
