package handlers

import (
	"authentication/util"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
)

func Welcome(w http.ResponseWriter, r *http.Request) render.Renderer {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return util.ErrUnauthorized(err)
		}

		return util.ErrBadRequest(err)
	}

	tknStr := c.Value

	claims := new(Claims)

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(t *jwt.Token) (interface{}, error) {
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

	return util.Response(fmt.Sprintf("Welcome %s!", claims.Username))
}
