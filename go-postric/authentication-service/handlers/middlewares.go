package handlers

import (
	"authentication/util"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
)

func AuthGuard() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			log.Printf("bearer token: %s", bearer)

			if bearer == "" {
				render.JSON(w, r, util.ErrResponse(http.StatusUnauthorized, errors.New("missing authorization header")))
				return
			}

			tknStr := strings.Replace(bearer, "Bearer ", "", 1)
			claims := new(Claims)

			tkn, err := jwt.ParseWithClaims(tknStr, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtSecretKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					render.JSON(w, r, util.ErrResponse(http.StatusUnauthorized, err))
					return
				}

				render.JSON(w, r, util.ErrResponse(http.StatusBadRequest, err))
				return
			}
			if !tkn.Valid {
				render.JSON(w, r, util.ErrResponse(http.StatusUnauthorized, err))
				return
			}

			ctxKey := util.ContextKey("claims")
			r = r.WithContext(context.WithValue(r.Context(), ctxKey, claims))
			next.ServeHTTP(w, r)
		})
	}
}
