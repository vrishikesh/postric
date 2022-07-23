package handlers

import (
	"authentication/util"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func Welcome(w http.ResponseWriter, r *http.Request) render.Renderer {
	ctxKey := util.ContextKey("claims")
	claims := r.Context().Value(ctxKey).(*Claims)

	return util.Response(http.StatusAccepted, fmt.Sprintf("Welcome %s!", claims.Username))
}
