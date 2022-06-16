package handlers

import (
	"authentication/util"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v4"
)

func Refresh(w http.ResponseWriter, r *http.Request) render.Renderer {
	ctxKey := util.ContextKey("claims")
	claims := r.Context().Value(ctxKey).(*Claims)

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return util.ErrResponse(http.StatusBadRequest, errors.New("refresh too early"))
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
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
